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

- Feature Name: (fill me in with a unique ident, `my_awesome_feature`)
- Start Date: (fill me in with today's date, YYYY-MM-DD)
- RFC PR: [rust-lang/rfcs#0000](https://github.com/rust-lang/rfcs/pull/0000)
- Rust Issue: [rust-lang/rust#0000](https://github.com/rust-lang/rust/issues/0000)

# Summary
[summary]: #summary

Horahora is a collaborative archival tool and media server. To enhance its collaborative functionality, we'd like to add "video rooms" with similar functionality to cytu.be or synctube: users would be able to join a room with an active owner, and would be able to watch synchronized videos with a small group of users. These rooms would contain a basic chat, and video queue.

# Motivation
[motivation]: #motivation

Why are we doing this? What use cases does it support? What is the expected outcome?

As mentioned, Horahora is a collaborative archival tool, and this feature enhances its role as a tool to collaboratively sort through archived content. This feature facilitates content sharing and content discovery for communities.

# Guide-level explanation
[guide-level-explanation]: #guide-level-explanation

Explain the proposal as if it was already included in the language and you were teaching it to another Rust programmer. That generally means:

1. a new "video room" page containing a list of rooms, which will have an option for each user to create their own room or join an existing room
2. a modified video page, consisting of:
    1. the recommended videos list replaced by a "watch queue", representing the order in which videos will be watched by the room's participants
    2. A basic text chat, pinned to e.g. the left of the screen
All other video page functionality (e.g. ratings, comments) will be left as is.

# Requirements
The solution MUST be:
1. relatively simple
2. Horizontally scalable
3. performant

# Reference-level explanation
[reference-level-explanation]: #reference-level-explanation

This feature will involve the creation of a new microservice: communityservice.

GRPC API will look like the following:

```service Communityservice {
    // Chat API
    rpc Chat(stream sendMessage) returns (respMessage) {}

    // Rooms API
    rpc CreateRoom(roomRequest) returns (roomResponse) {}
    rpc ListRoom(RoomListRequest) returns (roomListResponse) {}

    // Room video queue API
    rpc AddVideo(addVideoRequest) returns (Empty) {}
    rpc DeleteVideo(deleteVideoREquest) returns (Empty) {}
    rpc GetVideoQueue(getVideoQueueRequest) returns (videoQueueResponse) {}
}

message sendMessage {
    oneof Payload {
        TopicSubscription sub = 1;
        ChatMessage msg = 2;
    }
}

message respMessage {
    ChatMessage msg = 1;
}

message ChatMessage {
    string message = 1;
    string authorName = 2;
    timestamp timestamp = 3;
}

message TopicSubscription {
    uint32 room_id = 1;
    uint32 user_id = 2; // need this so we can determine whether to close the room if user leaves
}

message roomRequest {
    string title = 1;
    string ownerUID = 2;
}

message roomResponse {
    int roomID = 1;
}

message roomListRequest {
    uint32 page = 1;
}

message roomListResponse {
    repeated Room rooms = 1;
}

message Room {
    string title = 1;
    string ownerUsername = 2;
    int roomID = 3;
}```

Upon connecting to a room, the client will establish a websocket to front_api, and will first send the ID of the room that they've connected to. On receiving this, front_api will establish a GRPC stream using the Chat RPC to communityservice, first sending the ID of the room (received in the prior websocket message) to communityservice as a TopicSubscription message. On receiving the TopicSubscription message, communityservice will retrieve all messages for the room's history from the relevant Kafka topic (and start consuming from that topic), and send them to front_api, which will relay them to the client. The client will display all messages in order. 
On the client sending a message to front_api, front_api will send a ChatMessage into the corresponding GRPC stream. CommunityService, on this message's receipt, will immediately broadcast perform the following steps:
1. publish the message to the Kafka topic for the room
2. (in the consumption workers for other users connected) receive latest message for topic, send down GRPC stream, which will reach the client
If a new user joins the room, they will receive:
1. all messages starting from the beginning of the room
2. all new messages as they are produced

Rooms API is self explanatory, but will be stored in Postgres. Rooms schema:

CREATE TABLE rooms (
    id SERIAL primary key,
    name varchar(255),
    ownerUID uint32,
    creation_time timestamp
);

Rooms are closed when the owner leaves, tracked by TopicSubscription.UserID == ownerUID.

The video queue can be manipulated by the relevant GRPC functions (self-explanatory). 

# Drawbacks
[drawbacks]: #drawbacks

It bloats Horahora's setup, requiring Kafka and another microservice. 

# Rationale and alternatives
[rationale-and-alternatives]: #rationale-and-alternatives

We could alternatively have a separate page specific to rooms for watching videos, but it's best to leverage existing work. This solution keeps it simple (to the extent that we can) by slightly repurposing the existing video page for video rooms. Kafka is preferable to using Postgres as a performance consideration, and handles all of the message queueing, persistence, and retention for us, leaving communityservice stateless.

# Prior art
[prior-art]: #prior-art

I'm not aware of any uses of Kafka as a chat message store, but it definitely makes sense to use it as one, as it has all of the desirable properties: horizontally scalable, fast, control over in-flight (uncommitted) messages that can be queued, retention control, fault tolerance, etc.

# Unresolved questions
[unresolved-questions]: #unresolved-questions

None that come to mind.

# Future possibilities
[future-possibilities]: #future-possibilities

We could extend the chat in the future, e.g. to add emojis; this design complies with that vision.

Think about what the natural extension and evolution of your proposal would
be and how it would affect the language and project as a whole in a holistic
way. Try to use this section as a tool to more fully consider all possible
interactions with the project and language in your proposal.
Also consider how this all fits into the roadmap for the project
and of the relevant sub-team.

This is also a good place to "dump ideas", if they are out of scope for the
RFC you are writing but otherwise related.

If you have tried and cannot think of any future possibilities,
you may simply state that you cannot think of anything.

Note that having something written down in the future-possibilities section
is not a reason to accept the current or a future RFC; such notes should be
in the section on motivation or rationale in this or subsequent RFCs.
The section merely provides additional information.
