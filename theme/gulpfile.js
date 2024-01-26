const browsersync = require('browser-sync').create();
const cached = require('gulp-cached');
const cssnano = require('gulp-cssnano');
const del = require('del');
const fileinclude = require('gulp-file-include');
const gulp = require('gulp');
const gulpif = require('gulp-if');
const npmdist = require('gulp-npm-dist');
const replace = require('gulp-replace');
const uglify = require('gulp-uglify');
const useref = require('gulp-useref-plus');
const rename = require('gulp-rename');
const sourcemaps = require("gulp-sourcemaps");
const postcss = require('gulp-postcss');
const cleanCSS = require('gulp-clean-css');
const autoprefixer = require("gulp-autoprefixer");
const merge = require('merge-stream');

const paths = {
    config: {
        tailwindjs: "./tailwind.config.js",
    },
    base: {
        base: {
            dir: './'
        },
        node: {
            dir: './node_modules'
        },
        packageLock: {
            files: './package-lock.json'
        }
    },
    dist: {
        base: {
            dir: '../templates',
            files: '../templates/**/*'
        },
        baseassets: {
            dir: '../public',
            files: '../public/assets/**/*',
        },
        libs: {
            dir: '../public/assets/libs'
        },
        css: {
            dir: '../public/assets/css',
        },
        js: {
            dir: '../public/assets/js',
            files: '../public/assets/js/pages',
        },
        assetsembed: {
            dir: '../public',
            files: '../public/embed.go'
        },
        assetsrobots: {
            dir: '../public',
            files: '../public/robots.txt'
        },
        templatesembed: {
            dir: '../templates',
            files: '../templates/embed.go'
        }
    },
    src: {
        base: {
            dir: './src',
            files: './src/**/*'
        },
        html: {
            dir: './src',
            files: './src/**/*.html',
        },
        img: {
            dir: './src/assets/images',
            files: './src/assets/images/**/*',
        },
        js: {
            dir: './src/assets/js',
            pages: './src/assets/js/pages',
            files: './src/assets/js/pages/*.js',
            main: './src/assets/js/*.js',
        },
        partials: {
            dir: './src/partials',
            files: './src/partials/**/*'
        },
        css: {
            dir: './src/assets/css',
            files: './src/assets/css/**/*',
            icons: './src/assets/css/icons.css',
            main: './src/assets/css/*.css'
        },
        assetsembed: {
            dir: './src/assets',
            files: './src/assets/embed.go'
        },
        assetsrobots: {
            dir: './src/assets',
            files: './src/assets/robots.txt'
        },
        templatesembed: {
            dir: './src',
            files: './src/embed.go'
        },
        examples: {
            dir: './src/examples',
            files: './src/examples/**/*'
        },
        layouts: {
            dir: './src/layouts',
            files: './src/layouts/**/*'
        },
        pages: {
            dir: './src/pages',
            files: './src/pages/**/**/*'
        },
    }
};

gulp.task('browsersync', function (callback) {
    browsersync.init({
        port: 8000,
        server: {
            baseDir: [paths.dist.base.dir, paths.src.base.dir, paths.base.base.dir]
        },
    });
    callback();
});

gulp.task('browsersyncReload', function (callback) {
    browsersync.reload();
    callback();
});

gulp.task('watch', function () {
    gulp.watch([paths.src.css.files, paths.src.html.files,
        paths.src.layouts.files, paths.src.pages.files, paths.src.partials.files,
        "!" + paths.src.css.icons], gulp.series(['html', 'css']));
    gulp.watch([paths.src.css.icons], gulp.series('icons'));
    gulp.watch([paths.src.js.dir], gulp.series('js'));
    gulp.watch([paths.src.js.pages], gulp.series('jsPages'));
});

gulp.task('js', function () {
    return gulp
        .src(paths.src.js.main)
        .pipe(uglify())
        .pipe(gulp.dest(paths.dist.js.dir));
});

gulp.task('jsPages', function () {
    return gulp
        .src(paths.src.js.files)
        .pipe(uglify())
        .pipe(gulp.dest(paths.dist.js.files));
});

gulp.task('css', function () {

    // generate tailwind css
    return gulp
        .src(paths.src.css.main)
        .pipe(postcss())
        .pipe(gulp.dest(paths.dist.css.dir))
        .pipe(
            rename({
                suffix: ".min"
            })
        )
        .pipe(sourcemaps.init())
        .pipe(gulp.dest(paths.dist.css.dir));
});

gulp.task('icons', function () {

    return gulp
        .src(paths.src.css.icons)
        .pipe(postcss())
        .pipe(gulp.dest(paths.dist.css.dir))
        .pipe(
            rename({
                suffix: ".min"
            })
        )
        .pipe(sourcemaps.write("./"))
        .pipe(gulp.dest(paths.dist.css.dir));
});

gulp.task('fileinclude', function () {
    return gulp
        .src([
            paths.src.html.files,
            '!' + paths.dist.base.files,
            '!' + paths.src.partials.files
        ])
        .pipe(fileinclude({
            prefix: '@@',
            basepath: '@file',
            indent: true,
        }))
        .pipe(cached())
        .pipe(gulp.dest(paths.dist.base.dir));
});

gulp.task('clean:cache', function (callback) {
    cached.caches = {};
    callback();
});

gulp.task('clean:packageLock', function (callback) {
    del.sync(paths.base.packageLock.files);
    callback();
});

gulp.task('clean:dist', function (callback) {
    del.sync(paths.dist.base.dir, {
        force: true
    });
    callback();
});

gulp.task('copy:all', function () {
    let baseassets = gulp
        .src([
            paths.src.base.files,
            '!' + paths.src.partials.dir, '!' + paths.src.partials.files,
            '!' + paths.src.css.dir, '!' + paths.src.css.files,
            '!' + paths.src.js.dir, '!' + paths.src.js.files, '!' + paths.src.js.main,
            '!' + paths.src.html.files, '!' + paths.src.templatesembed.files,
            '!' + paths.src.assetsembed.files, '!' + paths.src.assetsrobots.files,
            '!' + paths.src.examples.dir, '!' + paths.src.examples.files,
            '!' + paths.src.layouts.dir, '!' + paths.src.layouts.files,
            '!' + paths.src.pages.dir, '!' + paths.src.pages.files,
        ])
        .pipe(gulp.dest(paths.dist.baseassets.dir));

    let assetsembed = gulp
        .src(paths.src.assetsembed.files)
        .pipe(gulp.dest(paths.dist.assetsembed.dir));

    let assetsrobots = gulp
        .src(paths.src.assetsrobots.files)
        .pipe(gulp.dest(paths.dist.assetsrobots.dir));

    return merge(baseassets, assetsembed, assetsrobots);
});

gulp.task('copy:libs', function () {
    return gulp
        .src(npmdist({
            replaceDefaultExcludes: true,
            excludes: [],
        }), { base: paths.base.node.dir })
        .pipe(rename(function (path) {
            path.dirname = path.dirname.replace(/\/dist/, '').replace(/\\dist/, '');
        }))
        .pipe(gulp.dest(paths.dist.libs.dir));
});

gulp.task('html:pages', function () {
    return gulp
        .src([
            paths.src.pages.files
        ])
        .pipe(cached())
        .pipe(gulp.dest(paths.dist.base.dir + '/pages'));
});

gulp.task('html:layouts', function () {
    return gulp
        .src([
            paths.src.layouts.files
        ])
        .pipe(cached())
        .pipe(gulp.dest(paths.dist.base.dir + '/layouts'));
});

gulp.task('html:partials', function () {
    return gulp
        .src([
            paths.src.partials.files
        ])
        .pipe(cached())
        .pipe(gulp.dest(paths.dist.base.dir + '/partials'));
});

gulp.task('html', function () {
    return gulp
        .src([
            paths.src.html.files, paths.src.templatesembed.files,
            '!' + paths.src.layouts.files, '!' + paths.src.layouts.dir,
            '!' + paths.src.partials.files, '!' + paths.src.partials.dir,
            '!' + paths.src.pages.files, '!' + paths.src.pages.dir,
            '!' + paths.src.examples.dir, '!' + paths.src.examples.files,
            '!' + paths.dist.baseassets.files, '!' + paths.dist.baseassets.dir,
        ])
        .pipe(useref())
        .pipe(cached())
        .pipe(gulpif('*.js', uglify()))
        .pipe(gulpif('*.css', cssnano({ svgo: false })))
        .pipe(gulp.dest(paths.dist.base.dir));
});

gulp.task('build', gulp.series(gulp.parallel('clean:packageLock', 'clean:dist', 'copy:all', 'copy:libs', 'css', 'icons', 'js', 'jsPages', 'html', 'html:layouts', 'html:partials', 'html:pages')));

gulp.task('default', gulp.series(gulp.parallel('clean:cache', 'clean:packageLock', 'clean:dist', 'copy:all', 'copy:libs', 'css', 'icons', 'js', 'jsPages', 'html', 'html:layouts', 'html:partials', 'html:pages'), gulp.parallel('watch')));