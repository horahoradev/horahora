var jqcomments = require('jquery-comments');

function loadComments(videoID, userID) {
    $.ajax({
        url: "/comments/" + videoID,
        dataType: 'json',
        success: function (response) {
            $('#comment-section').comments({
                profilePictureURL: '/static/images/placeholder1.jpg',
                getComments: function(success, error) {
                    success(response ?? {});
                },
                postComment: function(commentJSON, success, error) {
                    commentJSON['videoID'] = videoID;
                    commentJSON['userID'] = userID;

                    postData = {
                      video_id: videoID,
                      user_id: userID,
                      content: commentJSON.content,
                      parent: commentJSON.parent,
                    };

                    console.log(postData);

                    $.ajax({
                        type: 'post',
                        url: '/comments/',
                        data: postData,
                    });
                }
            });
        },
    });
}

// :(
window.loadComments = loadComments;