var jqcomments = require('jquery-comments');

function loadComments(videoID, userID) {
    $.ajax({
        url: "/comments/" + videoID,
        dataType: 'json',
        success: function (response) {
            $('#comment-section').comments({
                profilePictureURL: '/static/images/placeholder.jpg',
                getComments: function (success, error) {
                    success(response);
                },
                upvoteComment: function (commentJSON, success, error) {
                        postData = {
                            comment_id: commentJSON.id,
                            video_id: videoID,
                            user_id: userID,
                            user_has_upvoted: commentJSON.user_has_upvoted,
                        };

                        $.ajax({
                            type: 'post',
                            url: '/comment_upvotes/',
                            data: postData,
                        });
                    },
                postComment: function (commentJSON, success, error) {
                    commentJSON['videoID'] = videoID;
                    commentJSON['userID'] = userID;

                    postData = {
                        video_id: videoID,
                        user_id: userID,
                        content: commentJSON.content,
                        parent: commentJSON.parent,
                    };

                    $.ajax({
                        type: 'post',
                        url: '/comments/',
                        data: postData,
                    });
                },
            });
        }
    });
}

// :(
window.loadComments = loadComments;