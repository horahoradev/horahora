var jqcomments = require('jquery-comments');

function loadComments(videoID) {
    $.ajax({
        url: "/comments/" + videoID,
        dataType: 'json',
        success: function (response) {
            $('#comment-section').comments({
                profilePictureURL: '/static/images/placeholder1.jpg',
                getComments: function(success, error) {
                    success(response);
                }
            });
        }
    });
}

// :(
window.loadComments = loadComments;