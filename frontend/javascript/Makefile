build : video.js upload.js
	browserify video.js --s video -o ../assets/js/bundle.js && \
	browserify upload.js --s upload -o ../assets/js/upload.js && \
	cp node_modules/jquery-comments/css/* ../assets/css && \
	cp node_modules/dropzone/dist/dropzone.css ../assets/css