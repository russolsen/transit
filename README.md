# transit-go

Note that this readme file along with the rest of transit-go is very much a work-in-progress.

Transit is a data format and a set of libraries for conveying values between applications written in different languages. This library provides support for marshalling Transit data to/from Go.

* [Rationale](http://blog.cognitect.com/blog/2014/7/22/transit)
* [Specification](http://github.com/cognitect/transit-format)

This implementation's major.minor version number corresponds to the version of the Transit specification it supports.

Currently on the basic JSON format is implemented.
JSON-Verbose and MessagePack is **not** implemented yet. 

_NOTE: Transit is a work in progress and may evolve based on feedback. As a result, while Transit is a great option for transferring data between applications, it should not yet be used for storing data durably over time. This recommendation will change when the specification is complete._

## Usage

## Default Type Mapping

Note this is still TBD.

|Transit type|Write accepts|Read returns|
|------------|-------------|------------|
|null|nil|nil|
|string|string|string|
|bool|bool|bool|
|integer|int, int8, int16, int32, int64 and the unsigned variants|int64|
|decimal|float32, float64|float64|
|keyword|transit.Keyword|transit.Keyword|
|symbol|transit.Symbol|transit.Symbol|
|big decimal|big.|NForza.Transit.Numerics.BigRational|
|big integer|System.Numerics.BigInteger|System.Numerics.BigInteger|
|time|System.DateTime|System.DateTime|
|uri|System.Uri|System.Uri|
|uuid|System.Guid|System.Guid|
|char|System.Char|System.Char|
|array|T[], System.Collections.Generic.IList<>|System.Collections.Generic.IList<object>|
|list|System.Collections.Generic.IEnumerable<>|System.Collections.Generic.IEnumerable<object>|
|set|System.Collections.Generic.ISet<>|System.Collections.Generic.ISet<object>|
|map|System.Collections.Generic.IDictionary<,>|System.Collections.Generic.IDictionary<object, object>|
|link|NForza.Transit.ILink|NForza.Transit.ILink|
|ratio +|NForza.Transit.IRatio|NForza.Transit.IRatio|

\+ Extension type

## Layered Implementations

## Copyright and License
Copyright © 2016 Russ Olsen

This library is a Go port of the Java version created and maintained by Cognitect, therefore

Copyright © 2014 Cognitect

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at
http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.

This README file is based on the README from transit-csharp, therefore:

Copyright © 2014 NForza.

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at
http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.


