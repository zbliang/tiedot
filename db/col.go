/* Document collection. */
package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"loveoneanother.at/tiedot/file"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"
)

type IndexConf struct {
	FileName            string
	PerBucket, HashBits uint64
	IndexedPath         []string
}

type Config struct {
	Indexes []IndexConf
}

type Col struct {
	Data                                    *file.ColFile
	Config                                  *Config
	Dir, ConfigFileName, ConfBackupFileName string
	StrHT                                   map[string]*file.HashTable
	StrIC                                   map[string]*IndexConf
}

// Return string hash code.
func StrHash(thing interface{}) uint64 {
	// very similar to Java String.hashCode()
	// you must review (even rewrite) most collection test cases, if you change the hash algorithm
	str := fmt.Sprint(thing)
	hash := 0
	for _, c := range str {
		hash = int(c) + (hash << 6) + (hash << 16) - hash
	}
	return uint64(hash)
}

// Open a collection.
func OpenCol(dir string) (col *Col, err error) {
	if err = os.MkdirAll(dir, 0700); err != nil {
		return
	}
	col = &Col{ConfigFileName: path.Join(dir, "config"), ConfBackupFileName: path.Join(dir, "config.bak"), Dir: dir}
	// open data file
	if col.Data, err = file.OpenCol(path.Join(dir, "data")); err != nil {
		return
	}
	// make sure the config file exists
	tryOpen, err := os.OpenFile(col.ConfigFileName, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return
	} else if err = tryOpen.Close(); err != nil {
		return
	}
	col.LoadConf()
	return
}

// Copy existing config file content to backup config file.
func (col *Col) BackupAndSaveConf() error {
	oldConfig, err := ioutil.ReadFile(col.ConfigFileName)
	if err != nil {
		return err
	}
	if err = ioutil.WriteFile(col.ConfBackupFileName, []byte(oldConfig), 0600); err != nil {
		return err
	}
	if col.Config != nil {
		newConfig, err := json.Marshal(col.Config)
		if err != nil {
			return err
		}
		if err = ioutil.WriteFile(col.ConfigFileName, newConfig, 0600); err != nil {
			return err
		}
	}
	return nil
}

// (Re)load configuration to collection.
func (col *Col) LoadConf() error {
	// read index config
	config, err := ioutil.ReadFile(col.ConfigFileName)
	if err != nil {
		return err
	}
	if string(config) == "" {
		col.Config = &Config{}
	} else if err = json.Unmarshal(config, &col.Config); err != nil {
		return err
	}
	// open each index file
	col.StrHT = make(map[string]*file.HashTable)
	col.StrIC = make(map[string]*IndexConf)
	for i, index := range col.Config.Indexes {
		ht, err := file.OpenHash(path.Join(col.Dir, index.FileName), index.HashBits, index.PerBucket)
		if err != nil {
			return err
		}
		col.StrHT[strings.Join(index.IndexedPath, ",")] = ht
		col.StrIC[strings.Join(index.IndexedPath, ",")] = &col.Config.Indexes[i]
	}
	return nil
}

// Get inside the data structure, along the given path.
func GetIn(doc interface{}, path []string) (ret []interface{}) {
	thing := doc
	for _, seg := range path {
		if aMap, ok := thing.(map[string]interface{}); ok {
			thing = aMap[seg]
		} else {
			return nil
		}
	}
	switch thing.(type) {
	case []interface{}:
		return thing.([]interface{})
	default:
		return append(ret, thing)
	}
}

// Retrieve document data given its ID.
func (col *Col) Read(id uint64, doc interface{}) error {
	data := col.Data.Read(id)
	if data == nil {
		return errors.New(fmt.Sprintf("Document %d does not exist in %s", id, col.Dir))
	}
	if err := json.Unmarshal(data, &doc); err != nil {
		msg := fmt.Sprintf("Cannot parse document %d in %s to JSON, data: %s", id, col.Dir, string(data))
		log.Println(msg)       // for the srv
		return errors.New(msg) // for embedded usage
	}
	return nil
}

// Index the document on all indexes
func (col *Col) IndexDoc(id uint64, doc interface{}) {
	wg := new(sync.WaitGroup)
	wg.Add(len(col.StrIC))
	for k, v := range col.StrIC {
		go func(k string, v *IndexConf) {
			for _, thing := range GetIn(doc, v.IndexedPath) {
				if thing != nil {
					col.StrHT[k].Put(StrHash(thing), id)
				}
			}
			wg.Done()
		}(k, v)
	}
	wg.Wait()
}

// Remove the document from all indexes
func (col *Col) UnindexDoc(id uint64, doc interface{}) {
	wg := new(sync.WaitGroup)
	wg.Add(len(col.StrIC))
	for k, v := range col.StrIC {
		go func(k string, v *IndexConf) {
			for _, thing := range GetIn(doc, v.IndexedPath) {
				col.StrHT[k].Remove(StrHash(thing), 1, func(k, v uint64) bool {
					return v == id
				})
			}
			wg.Done()
		}(k, v)
	}
	wg.Wait()
}

// Insert a new document.
func (col *Col) Insert(doc interface{}) (id uint64, err error) {
	data, err := json.Marshal(doc)
	if err != nil {
		return
	}
	if id, err = col.Data.Insert(data); err != nil {
		return
	}
	col.IndexDoc(id, doc)
	return
}

// Insert a new document and immediately flush all buffers.
func (col *Col) DurableInsert(doc interface{}) (id uint64, err error) {
	id, err = col.Insert(doc)
	if err != nil {
		return
	}
	err = col.Flush()
	return
}

// Update a document, return its new ID.
func (col *Col) Update(id uint64, doc interface{}) (newID uint64, err error) {
	data, err := json.Marshal(doc)
	if err != nil {
		return
	}
	// read original document
	oldData := col.Data.Read(id)
	if oldData == nil {
		return id, errors.New(fmt.Sprintf("Document %d does not exist in %s", id, col.Dir))
	}
	var oldDoc interface{}
	if err = json.Unmarshal(oldData, &oldDoc); err == nil {
		col.UnindexDoc(id, oldDoc)
	} else {
		log.Printf("Original document %d in %s is corrupted, this update will attempt to overwrite it", id, col.Dir)
	}
	if newID, err = col.Data.Update(id, data); err != nil {
		return
	}
	col.IndexDoc(newID, doc)
	return
}

// Update a document and immediately flush all buffers.
func (col *Col) DurableUpdate(id uint64, doc interface{}) (newID uint64, err error) {
	newID, err = col.Update(id, doc)
	if err != nil {
		return
	}
	err = col.Flush()
	return
}

// Delete a document.
func (col *Col) Delete(id uint64) {
	var oldDoc interface{}
	err := col.Read(id, &oldDoc)
	if err != nil {
		return
	}
	col.Data.Delete(id)
	col.UnindexDoc(id, oldDoc)
}

// Delete a document and immediately flush all buffers.
func (col *Col) DurableDelete(id uint64) error {
	col.Delete(id)
	return col.Flush()
}

// Add an index.
func (col *Col) Index(path []string) error {
	joinedPath := strings.Join(path, ",")
	if _, found := col.StrHT[joinedPath]; found {
		return errors.New(fmt.Sprintf("Path %v is already indexed in collection %s", path, col.Dir))
	}
	newFileName := strings.Join(path, ",")
	if len(newFileName) > 100 {
		newFileName = newFileName[0:100]
	}
	// close all indexes
	for _, v := range col.StrHT {
		v.File.Close()
	}
	// save and reload config
	col.Config.Indexes = append(col.Config.Indexes, IndexConf{FileName: newFileName + strconv.Itoa(int(time.Now().UnixNano())), PerBucket: 200, HashBits: 14, IndexedPath: path})
	if err := col.BackupAndSaveConf(); err != nil {
		return err
	}
	if err := col.LoadConf(); err != nil {
		return err
	}
	// put all documents in the new index
	newIndex, ok := col.StrHT[strings.Join(path, ",")]
	if !ok {
		return errors.New(fmt.Sprintf("The new index %v in %s is gone??", path, col.Dir))
	}
	col.ForAll(func(id uint64, doc interface{}) bool {
		for _, thing := range GetIn(doc, path) {
			if thing != nil {
				newIndex.Put(StrHash(thing), id)
			}
		}
		return true
	})
	return nil
}

// Remove an index.
func (col *Col) Unindex(path []string) error {
	joinedPath := strings.Join(path, ",")
	if _, found := col.StrHT[joinedPath]; !found {
		return errors.New(fmt.Sprintf("Path %v was never indexed in collection %s", path, col.Dir))
	}
	// close all indexes
	for _, v := range col.StrHT {
		v.File.Close()
	}
	found := 0
	for i, index := range col.Config.Indexes {
		match := true
		for j, path := range path {
			if index.IndexedPath[j] != path {
				match = false
				break
			}
		}
		if match {
			found = i
			break
		}
	}
	// delete hash table file
	indexConf := col.Config.Indexes[found]
	indexHT := col.StrHT[strings.Join(indexConf.IndexedPath, ",")]
	indexHT.File.Close()
	if err := os.Remove(indexHT.File.Name); err != nil {
		return err
	}
	// remove it from config
	col.Config.Indexes = append(col.Config.Indexes[0:found], col.Config.Indexes[found+1:len(col.Config.Indexes)]...)
	if err := col.BackupAndSaveConf(); err != nil {
		return err
	}
	if err := col.LoadConf(); err != nil {
		return err
	}
	return nil
}

// Do fun for all documents (deserialized into generic interface).
func (col *Col) ForAll(fun func(id uint64, doc interface{}) bool) {
	col.Data.ForAll(func(id uint64, data []byte) bool {
		var parsed interface{}
		if err := json.Unmarshal(data, &parsed); err != nil {
			log.Printf("Cannot parse document '%v' in %s to JSON", data, col.Dir)
			return true
		} else {
			return fun(id, parsed)
		}
	})
}

// Do fun for all documents (deserialized into target interface). template must be a pointer to an initialized structure.
func (col *Col) DeserializeAll(template interface{}, fun func(id uint64) bool) {
	col.Data.ForAll(func(id uint64, data []byte) bool {
		if err := json.Unmarshal(data, template); err != nil {
			return true
		} else {
			return fun(id)
		}
	})
}

// Flush collection data files.
func (col *Col) Flush() error {
	if err := col.Data.File.Flush(); err != nil {
		log.Printf("Failed to flush %s, reason: %v", col.Data.File.Name, err)
		return err
	}
	for _, ht := range col.StrHT {
		if err := ht.File.Flush(); err != nil {
			log.Printf("Failed to flush %s, reason: %v", ht.File.Name, err)
			return err
		}
	}
	return nil
}

// Close a collection.
func (col *Col) Close() {
	if err := col.Data.File.Close(); err != nil {
		log.Printf("Failed to close %s, reason: %v", col.Data.File.Name, err)
	}
	for _, ht := range col.StrHT {
		if err := ht.File.Close(); err != nil {
			log.Printf("Failed to close %s, reason: %v", ht.File.Name, err)
		}
	}
}
