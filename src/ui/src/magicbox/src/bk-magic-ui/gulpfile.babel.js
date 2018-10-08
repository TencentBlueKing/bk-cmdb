/**
 * @file gulp 编译
 * @author ielgnaw <wuji0223@gmail.com>
 */

import gulp from 'gulp'
import sass from 'gulp-sass'
import replace from 'gulp-replace'
import gulpSequence from 'gulp-sequence'
import base64 from 'gulp-base64'
import autoprefixer from 'gulp-autoprefixer'
import cleanCSS from 'gulp-clean-css'
import rename from 'gulp-rename'
import sourceMap from 'gulp-sourcemaps'
import gzip from 'gulp-gzip'

gulp.task('compile-source', () => {
    return gulp.src(['./src/*.scss', '!./src/conf.scss'])
        // .pipe(base64(/*{maxImageSize: 300000}*/))
        .pipe(sourceMap.init())
        .pipe(sass())
        .pipe(autoprefixer({
            browsers: ['last 2 versions', 'ie > 8'],
            cascade: false
        }))
        .pipe(rename(path => {
            if (path.basename === 'bk') {
                path.basename = 'bk-magic-vue'
            }
        }))
        .pipe(sourceMap.write('.'))
        .pipe(gulp.dest('./lib/'))
})

gulp.task('compile-min', () => {
    return gulp.src(['./src/*.scss', '!./src/conf.scss'])
        // .pipe(base64(/*{maxImageSize: 300000}*/))
        .pipe(sourceMap.init())
        .pipe(sass({outputStyle: 'compressed'}))
        .pipe(autoprefixer({
            browsers: ['last 2 versions', 'ie > 8'],
            cascade: false
        }))
        // .pipe(cleanCSS())
        .pipe(rename(path => {
            if (path.basename === 'bk') {
                path.basename = 'bk-magic-vue'
            }
            path.extname = '.min' + path.extname
        }))
        .pipe(sourceMap.write('.'))
        .pipe(gulp.dest('./lib/'))
})

gulp.task('gzip-min', () => {
    return gulp.src(['./lib/bk-magic-vue.min.css'])
        .pipe(gzip({gzipOptions: {level: 9}}))
        .pipe(gulp.dest('./lib/'))
})

gulp.task('copy', () => {
    gulp.src('./src/fonts/**').pipe(gulp.dest('./lib/fonts'));
    gulp.src('./src/images/**').pipe(gulp.dest('./lib/images'));
});

// gulp.task('replace', () => {
//     return gulp.src('./dist/bk.css')
//         .pipe(replace('../image/', './image/'))
//         .pipe(replace('../font/', './font/'))
//         .pipe(gulp.dest('./dist/'));
// });

// gulp.task('replace4Source', () => {
//     let stream = gulp.src(['./src/**'])
//         .pipe(gulp.dest('./lib/'))
//         .pipe(replace('../image/', '/bkc-ui/lib/image/'))
//         .pipe(replace('../font/', '/bkc-ui/lib/font/'))
//         .pipe(gulp.dest('./lib/'));

//     return stream;
// });

gulp.task('base64', () => {
    return gulp.src('./lib/*.css')
        .pipe(base64(/*{maxImageSize: 300000}*/))
        // base64({
        //     baseDir: 'public',
        //     extensions: ['svg', 'png', /\.jpg#datauri$/i],
        //     exclude:    [/\.server\.(com|net)\/dynamic\//, '--live.jpg'],
        //     maxImageSize: 8*1024, // bytes
        //     debug: true
        // })
        .pipe(gulp.dest('./lib/'))
});

gulp.task('build', gulpSequence('copy', 'compile-source', 'compile-min', 'gzip-min' /*, 'base64'*/))
