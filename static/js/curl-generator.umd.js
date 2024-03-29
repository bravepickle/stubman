(function (global, factory) {
  typeof exports === 'object' && typeof module !== 'undefined' ? factory(exports) :
  typeof define === 'function' && define.amd ? define(['exports'], factory) :
  (global = typeof globalThis !== 'undefined' ? globalThis : global || self, factory(global["curl-generator"] = {}));
})(this, (function (exports) { 'use strict';

  // slash for connecting previous breakup line to current line for running cURL directly in Command Prompt
  var slash = " \\";
  var newLine = "\n";
  var newLinePadding = "  ";
  /**
   * @param {string} [method]
   * @returns {string}
   */
  var getCurlMethod = function (method) {
      var result = "";
      if (method) {
          var types = {
              GET: "-X GET",
              POST: "-X POST",
              PUT: "-X PUT",
              PATCH: "-X PATCH",
              DELETE: "-X DELETE",
          };
          result = types[method.toUpperCase()];
      }
      return slash + newLine + newLinePadding + result;
  };
  /**
   * @param {StringMap} headers
   * @returns {string}
   */
  var getCurlHeaders = function (headers) {
      var result = "";
      if (headers) {
          Object.keys(headers).map(function (val) {
              result += "" + slash + newLine + newLinePadding + "-H \"" + val + ": " + headers[val].replace(/(\\|")/g, "\\$1") + "\"";
          });
      }
      return result;
  };
  /**
   * @param {Object} body
   * @returns {string}
   */
  var getCurlBody = function (body) {
      var result = "";
      if (body) {
        var bodyVal;
        if (typeof body !== 'string' && !body instanceof String) {
            bodyVal = JSON.stringify(body)
        } else {
            bodyVal = body
        }

        result += '' + slash + newLine + newLinePadding + "--data-raw '" + bodyVal.replace(/(\\|')/g, "\\$1") + "'";
      }
      return result;
  };
  /**
   * Given the curl additional options, turn them into curl syntax
   * @param {CurlAdditionalOptions} [options]
   * @returns {string}
   */
  var getCurlOptions = function (options) {
      var result = "";
      if (options) {
          Object.keys(options).forEach(function (key) {
              var kebabKey = key.replace(/[A-Z]/g, function (letter) { return "-" + letter.toLowerCase(); });
              if (!options[key]) {
                  throw new Error("Invalid Curl option " + key);
              }
              else if (typeof options[key] === "boolean" && options[key]) {
                  // boolean option, we just add --opt
                  result += "--" + kebabKey + " ";
              }
              else if (typeof options[key] === "string") {
                  // string option, we have to add --opt=value
                  result += "--" + kebabKey + " " + options[key] + " ";
              }
          });
      }
      return result ? "" + slash + newLine + newLinePadding + result : result;
  };
  /**
   * @param {CurlRequest} params
   * @param {CurlAdditionalOptions} [options]
   * @returns {string}
   */
  var CurlGenerator = function (params, options) {
      var curlSnippet = "curl ";
      curlSnippet += params.url;
      curlSnippet += getCurlMethod(params.method);
      curlSnippet += getCurlHeaders(params.headers);
      curlSnippet += getCurlBody(params.body);
      curlSnippet += getCurlOptions(options);
      return curlSnippet.trim();
  };

  exports.CurlGenerator = CurlGenerator;

  Object.defineProperty(exports, '__esModule', { value: true });

}));
