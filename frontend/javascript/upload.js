/*
License

(The MIT License)

Copyright (c) 2012 Matias Meno <m@tias.me>
Logo & Website Design (c) 2015 "1910" www.weare1910.com

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

// I included the above license because I used and modified some of the code from
// https://www.dropzonejs.com/#usage and https://github.com/enyo/dropzone/wiki/Combine-normal-form-with-Dropzone
// when writing this function.

var dropzone = require('dropzone');

function setupUpload() {
    dropzone.autoDiscover = false;

    var uploadZone = new dropzone("#videoUpload", {
        url: "/upload",
        autoProcessQueue: false,
        uploadMultiple: false,
        maxFiles: 1,

        // The setting up of the dropzone
        init: function() {
            var myDropzone = this;

            // First change the button to actually tell Dropzone to process the queue.
            $("#upload").on("click", function(e) {
                // Make sure that the form isn't actually being sent.
                e.preventDefault();
                e.stopPropagation();
                myDropzone.processQueue();
            });
        },
    });

    uploadZone.on("sending", function(file, xhr, formData) {
        var title = $("#title").val();
        var description = $("#description").val();
        var tags = $("#tags").val(); // will likely need to change

        formData.append("title", title);
        formData.append("description", description);
        formData.append("tags", tags);

        console.log(formData);
    });
}

window.setupUpload = setupUpload;