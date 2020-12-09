"use strict";function _typeof(t){return(_typeof="function"==typeof Symbol&&"symbol"==typeof Symbol.iterator?function(t){return typeof t}:function(t){return t&&"function"==typeof Symbol&&t.constructor===Symbol&&t!==Symbol.prototype?"symbol":typeof t})(t)}!function(t){"function"==typeof define&&define.amd?define(["jquery"],t):"object"===("undefined"==typeof exports?"undefined":_typeof(exports))?t(require("jquery")):t(jQuery)}(function(l){var r=window.location,i="qor.sorting",t="enable."+i,o="change."+i,n="mousedown."+i,e="qor-sorting__highlight",p="qor-sorting__hover",f="tbody > tr";function a(t,o){this.$element=l(t),this.options=l.extend({},a.DEFAULTS,l.isPlainObject(o)&&o),this.$source=null,this.ascending=!1,this.orderType=0,this.startY=0,this.init()}return a.prototype={constructor:a,init:function(){var n,r=this.options,t=this.$element,o=t.find(f),i=0,e=0,s=0;l("body").addClass("qor-sorting"),t.find("tbody .qor-table__actions").append(a.TEMPLATE),o.each(function(t){var o=l(this).find(r.input).data("position");0<t&&(n<o?e++:e--),n=o,s=t}),e===s?i=1:-e===s&&(i=-1),this.$rows=o,this.orderType=i,this.bind()},bind:function(){var t=this.options;this.$element.on("click.qor.sorting",t.input,function(){return!1}).on(o,t.input,l.proxy(this.change,this)).on(n,t.toggle,l.proxy(this.mousedown,this)).on("mouseup.qor.sorting",l.proxy(this.mouseup,this)).on("dragstart.qor.sorting",f,l.proxy(this.dragstart,this)).on("dragend.qor.sorting",f,l.proxy(this.dragend,this)).on("dragover.qor.sorting",f,l.proxy(this.dragover,this)).on("drop.qor.sorting",f,l.proxy(this.drop,this))},unbind:function(){this.$element.off(o,this.change).off(n,this.mousedown)},change:function(t){var r,i=this.options,e=this.orderType,o=this.$rows,n=l(t.currentTarget),s=n.closest("tr"),a=s.parent(),u=n.data(),d=u.position,p=parseInt(n.val(),10),f=d<p;return t.stopPropagation(),o.each(function(){var t=l(this),o=t.find(i.input),n=o.data("position");n===p&&(r=t,f?1===e?r.after(s):-1===e&&r.before(s):1===e?r.before(s):-1===e&&r.after(s)),f?d<n&&n<=p&&o.data("position",--n).val(n):n<d&&p<=n&&o.data("position",++n).val(n)}),n.data("position",p),r||(f?1===e?a.append(s):-1===e&&a.prepend(s):1===e?a.prepend(s):-1===e&&a.append(s)),this.sort(s,{url:u.sortingUrl,from:d,to:p}),!1},mousedown:function(t){this.startY=t.pageY,l(t.currentTarget).closest("tr").prop("draggable",!0)},mouseup:function(){this.$element.find(f).prop("draggable",!1)},dragend:function(){l(f).removeClass(p),this.$element.find(f).prop("draggable",!1)},dragstart:function(t){var o=t.originalEvent,t=l(t.currentTarget);t.prop("draggable")&&o.dataTransfer&&(o.dataTransfer.effectAllowed="move",this.$source=t)},dragover:function(t){var o=this.$source;l(f).removeClass(p),l(t.currentTarget).prev("tr").addClass(p),o&&t.currentTarget!==this.$source[0]&&t.preventDefault()},drop:function(t){var o,n,r,i,e,s=this.options,a=this.orderType,u=t.pageY>this.startY,d=this.$source;l(f).removeClass(p),d&&t.currentTarget!==this.$source[0]&&(t.preventDefault(),n=l(t.currentTarget),t=(o=d.find(s.input)).data(),r=t.position,i=n.find(s.input).data("position"),e=r<i,this.$element.find(f).each(function(){var t=l(this).find(s.input),o=t.data("position");e?r<o&&o<=i&&t.data("position",--o).val(o):o<r&&i<=o&&t.data("position",++o).val(o)}),o.data("position",i).val(i),e?1===a||-1!==a&&u?n.after(d):n.before(d):1===a||-1!==a&&u?n.before(d):n.after(d),this.sort(d,{url:t.sortingUrl,from:r,to:i}))},sort:function(t,o){o.url&&(this.highlight(t),t.parent().after('<div class="qor-sorting__mask" style="position:absolute;width:100%;height:100%;top:0;left:0;z-index:10;background:#000;opacity:.4;"></div>'),l.ajax(o.url,{method:"post",data:{from:o.from,to:o.to},success:function(t,o,n){n.status},error:function(t,o,n){422===t.status?window.alert(t.responseText):window.alert([o,n].join(": ")),r.reload()}}).always(function(){l(".qor-sorting__mask").remove()}))},highlight:function(t){t.addClass(e),setTimeout(function(){t.removeClass(e)},2e3)},destroy:function(){this.unbind(),this.$element.removeData(i)}},a.DEFAULTS={toggle:!1,input:!1},a.TEMPLATE='<a class="qor-sorting__toggle"><i class="material-icons">drag_handle</i></a>',a.plugin=function(r){return this.each(function(){var t,o=l(this),n=o.data(i);if(!n){if(/destroy/.test(r))return;o.data(i,n=new a(this,r))}"string"==typeof r&&l.isFunction(t=n[r])&&t.apply(n)})},l(function(){var o,n;/sorting\=true/.test(r.search)&&(o=".qor-js-table",n={toggle:".qor-sorting__toggle",input:".qor-sorting__position"},l(document).on("disable.qor.sorting",function(t){a.plugin.call(l(o,t.target),"destroy")}).on(t,function(t){a.plugin.call(l(o,t.target),n)}).trigger("disable.qor.slideout").triggerHandler(t))}),a});