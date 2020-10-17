var jqcomments = require('jquery-comments');

function loadComments(videoID) {
    $.ajax({
        url: "/comments/" + videoID,
        dataType: 'json',
        success: function (response) {
            $('#comment-section').comments({
                profilePictureURL: 'https://app.viima.com/static/media/user_profiles/user-icon.png',
                getComments: function(success, error) {
                    success(response);
                }
            });
        }
    });
}

// :(
window.loadComments = loadComments;