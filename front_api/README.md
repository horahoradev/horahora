# REST API Documentation

### GET /home[?search=val][&page=x]
The query string val search is accepted, and will return videos whose title, tags, or description contains the search term. Inclusion and exclusion is supported, e.g. include1 include2 -exclude1

A query string val for the page number, starting at 1, is also accepted.

Response is of the form: {"PaginationData":{"PathsAndQueryStrings":["/home?page=1"],"Pages":[1],"CurrentPage":1},"Videos":[{"Title":"YOAKELAND","VideoID":1,"Views":6,"AuthorID":0,"AuthorName":"【旧】【旧】電ǂ鯨","ThumbnailLoc":"http://localhost:9000/otomads/7feaa38a-1e10-11ec-a6c3-0242ac1c0004.thumb","Rating":0}]}

For pagination data, the fields pages and PathsAndQueryStrings will always have the same length, and have corresponding values

Response is of the form: 

{"PaginationData":{"PathsAndQueryStrings":["/users/1?page=1"],"Pages":[1],"CurrentPage":1},"UserID":1,"Username":"【旧】【旧】電ǂ鯨","ProfilePictureURL":"/static/images/placeholder1.jpg","Videos":[{"Title":"YOAKELAND","VideoID":1,"Views":9,"AuthorID":0,"AuthorName":"【旧】【旧】電ǂ鯨","ThumbnailLoc":"http://localhost:9000/otomads/7feaa38a-1e10-11ec-a6c3-0242ac1c0004.thumb","Rating":0}]}

### GET /users/:id[?page=x] where id is the user id
A query string val for the page number, starting at 1, is accepted.

Response is of the form:

{"PaginationData":{"PathsAndQueryStrings":["/users/1?page=1"],"Pages":[1],"CurrentPage":1},"UserID":1,"Username":"【旧】【旧】電ǂ鯨","ProfilePictureURL":"/static/images/placeholder1.jpg","Videos":[{"Title":"YOAKELAND","VideoID":1,"Views":11,"AuthorID":0,"AuthorName":"【旧】【旧】電ǂ鯨","ThumbnailLoc":"http://localhost:9000/otomads/7feaa38a-1e10-11ec-a6c3-0242ac1c0004.thumb","Rating":0}]}

For pagination data, the fields pages and PathsAndQueryStrings will always have the same length, and have corresponding values

### GET /videos/:id where id is the video id
Response is of the form:

{"Title":"コダック","MPDLoc":"http://localhost:9000/otomads/207f773c-1e23-11ec-a6c3-0242ac1c0004.mpd","Views":2,"Rating":0,"VideoID":5,"AuthorID":5,"Username":"たっぴ","UserDescription":"","VideoDescription":"YouTube　\u003ca href=\"https://youtu.be/kP_lYd9D2to\" target=\"_blank\" rel=\"noopener nofollow\"\u003ehttps://youtu.be/kP_lYd9D2to\u003c/a\u003e","UserSubscribers":0,"ProfilePicture":"/static/images/placeholder1.jpg","UploadDate":"2021-09-25T17:07:56.400857Z","Comments":null,"Tags":null}

AuthorID, userDescription, and userSubscribers all have no meaning as of yet.

### GET /comments/:id , where :id is the video ID
Response is of this form: 

[{"id":1,"created":"2021-09-25T16:46:53.141031Z","content":"test","fullname":"admin","profile_picture_url":"/static/images/placeholder.png","upvote_count":0,"user_has_upvoted":false}]

### GET /archiverequests
Requires authentication

route: GET /archiverequests

Response is of this form:

{"ArchivalRequests": [{"url":"https://www.youtube.com/watch?v=8DXqneHHzA8"}]}

### POST /login
Accepts form-encoded values: username, password

response: 200 if ok, and sets a cookie

### POST /register
Accepts form-encoded values username, password, and email

response: 200 if ok, and sets a cookie

### POST /logout
Accepts no parameters

response: 200 if ok

### POST /archiverequests
Requires authentication

Accepts form-encoded value URL, which is the url to be archived

response: 200 if ok

### POST /rate/:id where :id is the video id
Accepts query parameter "rating" (float) 

Requires authentication

response: 200 if ok

### POST /approve/:id where :id is the video id
Requires authentication

Allows the user, if sufficiently high rank, to approve of a video and allow it to be shown to regular users.

Response: 200 if okay

### POST /comments/
Requires authentication

Accepts form-encoded values: video_id, content (content of comment), and parent (parent comment id if a reply)

response: 200 if ok

### POST /comment_upvotes/
Requires authentication

Accepts form-encoded value comment_id, which is the url to be archived

response: 200 if ok

