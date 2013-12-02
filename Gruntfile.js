module.exports = function(grunt) {

	grunt.initConfig({
		pkg: grunt.file.readJSON("package.json"),

		build_dir: "build",
		compile_dir: "bin",

		app_files: {
			js: ['web/src/**/*.js', '!web/src/assets/**/*.js'],
			atpl: ['web/src/app/**/*.tpl.html'],
			html: ["web/src/index.html"],
			css: ["web/css/**/*.css"]
		},

		lib_files: {
			js: [
				"web/lib/angular/angular.js",
				"web/lib/angular-ui-router/release/angular-ui-router.js",
				"web/lib/angular-ui-utils/modules/route/route.js"
			]
		},

		clean: [
			"<%= build_dir %>",
			"<%= compile_dir %>"
		],

		copy: {
			build_appjs: {
				files: [
					{
						src: ["<%= app_files.js %>", "<%= app_files.css %>"],
						dest: "<%= build_dir %>",
						cwd: ".",
						expand: true,
						flatten: true,
					}
				]
			},
			build_libjs: {
				files: [
					{
						src: ["<%= lib_files.js %>"],
						dest: "<%= build_dir %>",
						cwd: ".",
						expand: true,
						flatten: true,
					}
				]
			}
		},

	    html2js: {
	      /**
	       * These are the templates from `src/app`.
	       */
	      app: {
	        options: {
	          base: 'web/src/app'
	        },
	        src: [ '<%= app_files.atpl %>' ],
	        dest: '<%= build_dir %>/templates-app.js'
	      }

	    },		

		index: {
			build: {
				dir: "<%= build_dir %>",
				lib: "<%= lib_files.js %>",
				app: ["web/src/app.js", "web/src/home/home.js", "web/src/movies/movies.js", "web/src/common/core.js"],
				tpl: "<%= html2js.app.dest %>"
			}
		}
	});

	grunt.loadNpmTasks("grunt-contrib-clean");
	grunt.loadNpmTasks("grunt-contrib-copy");
	grunt.loadNpmTasks("grunt-html2js");

	grunt.registerMultiTask("index", "Process index.html template", function() {
		var data = [];
		data = data.concat(this.data.lib, this.data.app, this.data.tpl);

		var jsFiles = data.map( function(file) {
			return file.replace(/^.*[\\\/]/, '');
		});

		grunt.file.copy("web/src/index.html", this.data.dir + "/index.html", {
			process: function(contents, path) {
				return grunt.template.process(contents, {
					data: {
						scripts: jsFiles,
						version: grunt.config("pkg.version")
					}
				})
			}
		});
	});

	grunt.registerTask("default", ["clean", "html2js", "copy:build_appjs", "copy:build_libjs", "index:build"])
}