/******/ var __webpack_modules__ = ({

/***/ 886:
/***/ ((module) => {

function _extends() {
  module.exports = _extends = Object.assign || function (target) {
    for (var i = 1; i < arguments.length; i++) {
      var source = arguments[i];

      for (var key in source) {
        if (Object.prototype.hasOwnProperty.call(source, key)) {
          target[key] = source[key];
        }
      }
    }

    return target;
  };

  module.exports["default"] = module.exports, module.exports.__esModule = true;
  return _extends.apply(this, arguments);
}

module.exports = _extends;
module.exports["default"] = module.exports, module.exports.__esModule = true;

/***/ }),

/***/ 56:
/***/ ((module, __unused_webpack_exports, __nccwpck_require__) => {

var objectWithoutPropertiesLoose = __nccwpck_require__(503);

function _objectWithoutProperties(source, excluded) {
  if (source == null) return {};
  var target = objectWithoutPropertiesLoose(source, excluded);
  var key, i;

  if (Object.getOwnPropertySymbols) {
    var sourceSymbolKeys = Object.getOwnPropertySymbols(source);

    for (i = 0; i < sourceSymbolKeys.length; i++) {
      key = sourceSymbolKeys[i];
      if (excluded.indexOf(key) >= 0) continue;
      if (!Object.prototype.propertyIsEnumerable.call(source, key)) continue;
      target[key] = source[key];
    }
  }

  return target;
}

module.exports = _objectWithoutProperties;
module.exports["default"] = module.exports, module.exports.__esModule = true;

/***/ }),

/***/ 503:
/***/ ((module) => {

function _objectWithoutPropertiesLoose(source, excluded) {
  if (source == null) return {};
  var target = {};
  var sourceKeys = Object.keys(source);
  var key, i;

  for (i = 0; i < sourceKeys.length; i++) {
    key = sourceKeys[i];
    if (excluded.indexOf(key) >= 0) continue;
    target[key] = source[key];
  }

  return target;
}

module.exports = _objectWithoutPropertiesLoose;
module.exports["default"] = module.exports, module.exports.__esModule = true;

/***/ })

/******/ });
/************************************************************************/
/******/ // The module cache
/******/ var __webpack_module_cache__ = {};
/******/ 
/******/ // The require function
/******/ function __nccwpck_require__(moduleId) {
/******/ 	// Check if module is in cache
/******/ 	var cachedModule = __webpack_module_cache__[moduleId];
/******/ 	if (cachedModule !== undefined) {
/******/ 		return cachedModule.exports;
/******/ 	}
/******/ 	// Create a new module (and put it into the cache)
/******/ 	var module = __webpack_module_cache__[moduleId] = {
/******/ 		// no module.id needed
/******/ 		// no module.loaded needed
/******/ 		exports: {}
/******/ 	};
/******/ 
/******/ 	// Execute the module function
/******/ 	var threw = true;
/******/ 	try {
/******/ 		__webpack_modules__[moduleId](module, module.exports, __nccwpck_require__);
/******/ 		threw = false;
/******/ 	} finally {
/******/ 		if(threw) delete __webpack_module_cache__[moduleId];
/******/ 	}
/******/ 
/******/ 	// Return the exports of the module
/******/ 	return module.exports;
/******/ }
/******/ 
/************************************************************************/
/******/ /* webpack/runtime/define property getters */
/******/ (() => {
/******/ 	// define getter functions for harmony exports
/******/ 	__nccwpck_require__.d = (exports, definition) => {
/******/ 		for(var key in definition) {
/******/ 			if(__nccwpck_require__.o(definition, key) && !__nccwpck_require__.o(exports, key)) {
/******/ 				Object.defineProperty(exports, key, { enumerable: true, get: definition[key] });
/******/ 			}
/******/ 		}
/******/ 	};
/******/ })();
/******/ 
/******/ /* webpack/runtime/hasOwnProperty shorthand */
/******/ (() => {
/******/ 	__nccwpck_require__.o = (obj, prop) => (Object.prototype.hasOwnProperty.call(obj, prop))
/******/ })();
/******/ 
/******/ /* webpack/runtime/compat */
/******/ 
/******/ if (typeof __nccwpck_require__ !== 'undefined') __nccwpck_require__.ab = new URL('.', import.meta.url).pathname.slice(import.meta.url.match(/^file:\/\/\/\w:/) ? 1 : 0, -1) + "/";
/******/ 
/************************************************************************/
var __webpack_exports__ = {};
// This entry need to be wrapped in an IIFE because it need to be isolated against other modules in the chunk.
(() => {

// EXPORTS
__nccwpck_require__.d(__webpack_exports__, {
  "B": () => (/* reexport */ createHighlightComponent),
  "X": () => (/* reexport */ createHitsComponent),
  "cx": () => (/* reexport */ cx)
});

// EXTERNAL MODULE: ../@babel/runtime/helpers/extends.js
var helpers_extends = __nccwpck_require__(886);
// EXTERNAL MODULE: ../@babel/runtime/helpers/objectWithoutProperties.js
var objectWithoutProperties = __nccwpck_require__(56);
;// CONCATENATED MODULE: ./dist/es/lib/cx.js
function cx() {
  for (var _len = arguments.length, classNames = new Array(_len), _key = 0; _key < _len; _key++) {
    classNames[_key] = arguments[_key];
  }
  return classNames.reduce(function (acc, className) {
    if (Array.isArray(className)) {
      return acc.concat(className);
    }
    return acc.concat([className]);
  }, []).filter(Boolean).join(' ');
}
;// CONCATENATED MODULE: ./dist/es/components/Highlight.js


var _excluded = ["parts", "highlightedTagName", "nonHighlightedTagName", "separator", "className", "classNames"];

function createHighlightPartComponent(_ref) {
  var createElement = _ref.createElement;
  return function HighlightPart(_ref2) {
    var classNames = _ref2.classNames,
      children = _ref2.children,
      highlightedTagName = _ref2.highlightedTagName,
      isHighlighted = _ref2.isHighlighted,
      nonHighlightedTagName = _ref2.nonHighlightedTagName;
    var TagName = isHighlighted ? highlightedTagName : nonHighlightedTagName;
    return createElement(TagName, {
      className: isHighlighted ? classNames.highlighted : classNames.nonHighlighted
    }, children);
  };
}
function createHighlightComponent(_ref3) {
  var createElement = _ref3.createElement,
    Fragment = _ref3.Fragment;
  var HighlightPart = createHighlightPartComponent({
    createElement: createElement,
    Fragment: Fragment
  });
  return function Highlight(userProps) {
    var parts = userProps.parts,
      _userProps$highlighte = userProps.highlightedTagName,
      highlightedTagName = _userProps$highlighte === void 0 ? 'mark' : _userProps$highlighte,
      _userProps$nonHighlig = userProps.nonHighlightedTagName,
      nonHighlightedTagName = _userProps$nonHighlig === void 0 ? 'span' : _userProps$nonHighlig,
      _userProps$separator = userProps.separator,
      separator = _userProps$separator === void 0 ? ', ' : _userProps$separator,
      className = userProps.className,
      _userProps$classNames = userProps.classNames,
      classNames = _userProps$classNames === void 0 ? {} : _userProps$classNames,
      props = objectWithoutProperties(userProps, _excluded);
    return createElement("span", helpers_extends({}, props, {
      className: cx(classNames.root, className)
    }), parts.map(function (part, partIndex) {
      var isLastPart = partIndex === parts.length - 1;
      return createElement(Fragment, {
        key: partIndex
      }, part.map(function (subPart, subPartIndex) {
        return createElement(HighlightPart, {
          key: subPartIndex,
          classNames: classNames,
          highlightedTagName: highlightedTagName,
          nonHighlightedTagName: nonHighlightedTagName,
          isHighlighted: subPart.isHighlighted
        }, subPart.value);
      }), !isLastPart && createElement("span", {
        className: classNames.separator
      }, separator));
    }));
  };
}
;// CONCATENATED MODULE: ./dist/es/components/Hits.js


var Hits_excluded = ["classNames", "hits", "itemComponent", "sendEvent", "emptyComponent"];


// Should be imported from a shared package in the future

function createHitsComponent(_ref) {
  var createElement = _ref.createElement;
  return function Hits(userProps) {
    var _userProps$classNames = userProps.classNames,
      classNames = _userProps$classNames === void 0 ? {} : _userProps$classNames,
      hits = userProps.hits,
      ItemComponent = userProps.itemComponent,
      sendEvent = userProps.sendEvent,
      EmptyComponent = userProps.emptyComponent,
      props = objectWithoutProperties(userProps, Hits_excluded);
    if (hits.length === 0 && EmptyComponent) {
      return createElement(EmptyComponent, {
        className: cx('ais-Hits', classNames.root, cx('ais-Hits--empty', classNames.emptyRoot), props.className)
      });
    }
    return createElement("div", helpers_extends({}, props, {
      className: cx('ais-Hits', classNames.root, hits.length === 0 && cx('ais-Hits--empty', classNames.emptyRoot), props.className)
    }), createElement("ol", {
      className: cx('ais-Hits-list', classNames.list)
    }, hits.map(function (hit, index) {
      return createElement(ItemComponent, {
        key: hit.objectID,
        hit: hit,
        index: index,
        className: cx('ais-Hits-item', classNames.item),
        onClick: function onClick() {
          sendEvent('click:internal', hit, 'Hit Clicked');
        },
        onAuxClick: function onAuxClick() {
          sendEvent('click:internal', hit, 'Hit Clicked');
        }
      });
    })));
  };
}
;// CONCATENATED MODULE: ./dist/es/components/index.js


;// CONCATENATED MODULE: ./dist/es/lib/index.js

;// CONCATENATED MODULE: ./dist/es/index.js



})();

var __webpack_exports__createHighlightComponent = __webpack_exports__.B;
var __webpack_exports__createHitsComponent = __webpack_exports__.X;
var __webpack_exports__cx = __webpack_exports__.cx;
export { __webpack_exports__createHighlightComponent as createHighlightComponent, __webpack_exports__createHitsComponent as createHitsComponent, __webpack_exports__cx as cx };
