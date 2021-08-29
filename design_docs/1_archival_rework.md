Permission is hereby granted, free of charge, to any
person obtaining a copy of this software and associated
documentation files (the "Software"), to deal in the
Software without restriction, including without
limitation the rights to use, copy, modify, merge,
publish, distribute, sublicense, and/or sell copies of
the Software, and to permit persons to whom the Software
is furnished to do so, subject to the following
conditions:

The above copyright notice and this permission notice
shall be included in all copies or substantial portions
of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF
ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED
TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A
PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT
SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY
CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR
IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER
DEALINGS IN THE SOFTWARE.

- Feature Name: Archival rework
- Start Date: 8-17-2021
- RFC PR: 
- Rust Issue: 

# Summary
[summary]: #summary

One paragraph explanation of the feature.
Archive requests should be reworked to be website and structure-agnostic. The code should be simplified and generalized, allowing the project to support all websites and content forms that yt-dlp supports.

# Motivation
[motivation]: #motivation

Currently, Horahora has very specific logic regarding content types and supported websites; this logic is even baked into protobuf definitions, making it very laborious to expand site or content support.


# Guide-level explanation
[guide-level-explanation]: #guide-level-explanation

The workflow goes like this:
1. When a user requests a piece of content, it's immediately created as an archival request. Upon creation, we assess whether the URL belongs to a special class of yt-dlp extractors which return results in ascending order. If it is, we mark it as such.
2. Once a piece of content is created, it's acted upon by scheduler's sync manager workers. This worker will fetch an archival request whose backoff period has expired, and download all videos for it. If the archival request is "tailable" (extractor returns results in ascending order), we pass in some special flags to youtube-dl indicating that it can stop once it reaches the last known result.
3. The scheduler's scheduler workers will fetch the latest video in a given category randomly with probability according to its number of subscribers.
4. When the scheduler's downloader worker uploads the video, videoservice will store the original website string and video ID. The website is extracted from the domain name.

# Reference-level explanation
[reference-level-explanation]: #reference-level-explanation

The guide-level explanation explains this in sufficient detail.


# Drawbacks
[drawbacks]: #drawbacks

Why should we *not* do this?

THe logic pertaining to tailing on a per-extractor basis is hard to maintain and a little fragile. Otherwise, though, no drawbacks.

# Rationale and alternatives
[rationale-and-alternatives]: #rationale-and-alternatives

- Why is this design the best in the space of possible designs?
It's generic, and makes it much much easier to add support for new sites and video organizations.

- What other designs have been considered and what is the rationale for not choosing them?
We could ditch the storage of videos to queue entirely and simply let yt-dlp run for e.g. one video. The problem there would be that it would require yt-dlp to make additional requests to reach, for example, page 5.

- What is the impact of not doing this?
Horahora languishes in obscurity forever, with limited site and content support.

# Prior art
[prior-art]: #prior-art

Nothing relevant to discuss here.

# Unresolved questions
[unresolved-questions]: #unresolved-questions

- What parts of the design do you expect to resolve through the RFC process before this gets merged?
All, more or less (see above description)

- What parts of the design do you expect to resolve through the implementation of this feature before stabilization?
None

- What related issues do you consider out of scope for this RFC that could be addressed in the future independently of the solution that comes out of this RFC?
None

# Future possibilities
[future-possibilities]: #future-possibilities

We could improve the website identification for websites with multiple domains. I don't have any great ideas here yet, maybe we could use WHOIS.
