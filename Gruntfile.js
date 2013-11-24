module.exports = function(grunt) {

	grunt.initConfig({
		pkg: grunt.file.readJSON("package.json"),

		build_dir: "build",
		compile_dir: "bin",

		app_files: {
			js: ["web/src/**/*.js", "!web/src/assets/**/*.js"],
			atpl: ["web/src/app/**/*.tpl.html"],
			html: ["web/src/index.html"]
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
						src: ["<%= app_files.js %>"],
						dest: "<%= build_dir %>",
						cwd: ".",
						expand: true,
						flatten: true
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
						flatten: true
					}
				]
			}
		},

		index: {
			build: {
				dir: "<%= build_dir %>",
				js: [
					"<%= lib_files %>"
				]
			}
		}
	});

	grunt.loadNpmTasks("grunt-contrib-clean");
	grunt.loadNpmTasks("grunt-contrib-copy");

	grunt.registerMultiTask("index", "Process index.html template", function() {
		grunt.file.copy("web/src/index.html", this.data.dir + "/index.html", {
			process: function(contents, path) {
				return grunt.template.process(contents, {
					data: {
						scripts: grunt.config("this.data.js"),
						version: grunt.config("pkg.version")
					}
				})
			}
		});
	});

	grunt.registerTask("default", ["clean", "copy:build_appjs", "copy:build_libjs", "index:build"])
}