> To use `go get`, please use URL `loveoneanother.at/tiedot`. See INSTALL for more details.

<img src="http://golang.org/doc/gopher/frontpage.png" alt="Golang logo" align="right"/>

### tiedot - Your NoSQL database powered by Golang

tiedot is a document database that uses __JSON__ for documents and queries; it can be __embedded__ into your program, or run a stand-alone server using __HTTP__ for an API.

### Feature Highlights

- Designed for both embedded usage and standalone service.
- Fault-tolerant data structures that put safety of your data *first*.
- Very scalable on SMP computers.
- Use JSON syntax to build powerful queries.
- Support both \*nix and Windows operating systems.

### High Performance!

tiedot scales very well on SMP computers. The following performance results are collected from three machines types, using tiedot built-in benchmark:

(Operations per second)
<table>
<tr>
  <th>Processor</th>
  <th>Insert</th>
  <th>Read</th>
  <th>Query</th>
  <th>Update</th>
  <th>Delete</th>
  <th>Mix*</th>
  <th>Machine Type</th>
</tr>
<tr>
  <td>Mobile Intel Core i7 (2nd Gen)</td>
  <td>168k</td>
  <td>300k</td>
  <td>66k</td>
  <td>58k</td>
  <td>138k</td>
  <td>145k</td>
  <td>A 3 years old laptop</td>
</tr>
<tr>
  <td>Desktop Intel Core 2 Duo</td>
  <td>127k</td>
  <td>248k</td>
  <td>53k</td>
  <td>56k</td>
  <td>125k</td>
  <td>121k</td>
  <td>A 5 years old workstation</td>
</tr>
<tr>
  <td>Amazon EC2 c1.medium</td>
  <td>67k</td>
  <td>120k</td>
  <td>29k</td>
  <td>32k</td>
  <td>100k</td>
  <td>63k</td>
  <td>A medium range instance</td>
</tr>
</table>

Mix\* runs simultaneous insert/read/update/delete/query operations. See [Performance tuning and benchmarks] for more details

### References

- [Quick Start Guide: tiedot in 10 minutes]
- [API V1 reference]
- [API V2 reference]
- [Embedded usage]
- [Data structures]
- [Query processor and index]
- [Concurrency and networking]
- [Performance tuning and benchmarks]
- [Limitations]

### Version History

See [Version History] for detailed change logs, known issues, etc.

<table>
<tr>
  <th>Branch</th>
  <th>Release Version</th>
  <th>API Support</th>
  <th>Release Date</th>
  <th>Highlights</th>
</tr>
<tr>
  <td>alpha</td>
  <td>Alpha</td>
  <td>V1 only</td>
  <td>2013-06-28</td>
  <td>First release</td>
</tr>
<tr>
  <td>beta</td>
  <td>Beta</td>
  <td>V1 only</td>
  <td>2013-07-12</td>
  <td>Platform support and data durability improvements.</td>
</tr>
<tr>
  <td>1.0</td>
  <td>1.0</td>
  <td>V1 and V2</td>
  <td>2013-09-21</td>
  <td>Query performance/syntax, and documentation improvements.</td>
</tr>
</table>

### Contact and License

Future development plans are tracked in [Issues] section.

Please contact [Howard] for feedback /questions. I would love to hear from you! Please also check out my [Twitter] and [blog].

tiedot uses 2-clause BSD license:
<pre>
Copyright (c) 2013, Howard Guo
All rights reserved.

Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:
- Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
- Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
</pre>


### Project Background

__Is tiedot "yet another NoSQL database"?__

There are probably as many NoSQL database as there are Linux distributions.

tiedot is not as powerful (yet) - and does not intend to compete with mainstream NoSQL database engines such as CouchDB or Cassandra. However, tiedot performs reasonably well given its small size (around 3k LOC); and due to its simplicity, its performance may come close or even exceed those large brand NoSQL solutions (under certain workloads).

__What is the motive behind this project?__

Golang (Go) is a fascinating language - very easy to use, scalable and reasonably stable. I am very passionate about document database technologies (check out my other GitHub projects!) and enjoy seeing my code scaling well on SMP machines. This is my Golang exercise.

__Why the name "tiedot"?__

"tiedot" is a Finnish word standing for "data". I enjoy learning (natural and computer) languages, also enjoy listening to music in many languages. "Tiedot" sounds cute, doesn't it?

[Quick Start Guide: tiedot in 10 minutes]: https://github.com/HouzuoGuo/tiedot/wiki/Tutorial
[API V1 reference]: https://github.com/HouzuoGuo/tiedot/wiki/API-V1-Reference
[API V2 reference]: https://github.com/HouzuoGuo/tiedot/wiki/API-V2-Reference
[Version History]: https://github.com/HouzuoGuo/tiedot/wiki/Version-History
[Embedded usage]: https://github.com/HouzuoGuo/tiedot/wiki/Embedded-Usage
[Data structures]: https://github.com/HouzuoGuo/tiedot/wiki/Data-structures
[Query processor and index]: https://github.com/HouzuoGuo/tiedot/wiki/Query-processor-and-index
[Concurrency and networking]: https://github.com/HouzuoGuo/tiedot/wiki/Concurrency-and-networking
[Performance tuning and benchmarks]: https://github.com/HouzuoGuo/tiedot/wiki/Performance-tuning-and-benchmarks
[Limitations]: https://github.com/HouzuoGuo/tiedot/wiki/Limitations
[Howard]: mailto:guohouzuo@gmail.com
[Twitter]: https://twitter.com/hzguo
[blog]: http://allstarnix.blogspot.com.au
[Issues]: https://github.com/HouzuoGuo/tiedot/issues