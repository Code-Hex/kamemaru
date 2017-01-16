var gulp       = require('gulp')
var browserify = require('browserify')
var source     = require('vinyl-source-stream')

gulp.task('build', function() {
  browserify({
    'entries': ['./client/main.js']
  }).bundle()
    .pipe(source('bundle.js'))
    .pipe(gulp.dest('./assets/js'))
})

gulp.task('default', ['build'])
