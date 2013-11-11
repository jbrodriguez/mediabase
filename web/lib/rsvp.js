(function(e){"use strict";function v(e,n){t.async(function(){e.trigger("promise:resolved",{detail:n}),e.isResolved=!0,e.resolvedValue=n})}function m(e,n){t.async(function(){e.trigger("promise:failed",{detail:n}),e.isRejected=!0,e.rejectedValue=n})}function g(e){var t,n=[],r=new h,i=e.length;i===0&&r.resolve([]);var s=function(e){return function(t){o(e,t)}},o=function(e,t){n[e]=t,--i===0&&r.resolve(n)},u=function(e){r.reject(e)};for(t=0;t<i;t++)e[t].then(s(t),u);return r}function y(e,n){t[e]=n}var t={},n=typeof window!="undefined"?window:{},r=n.MutationObserver||n.WebKitMutationObserver,i;if(typeof process!="undefined"&&{}.toString.call(process)==="[object process]")t.async=function(e,t){process.nextTick(function(){e.call(t)})};else if(r){var s=[],o=new r(function(){var e=s.slice();s=[],e.forEach(function(e){var t=e[0],n=e[1];t.call(n)})}),u=document.createElement("div");o.observe(u,{attributes:!0}),window.addEventListener("unload",function(){o.disconnect(),o=null}),t.async=function(e,t){s.push([e,t]),u.setAttribute("drainQueue","drainQueue")}}else t.async=function(e,t){setTimeout(function(){e.call(t)},1)};var a=function(e,t){this.type=e;for(var n in t){if(!t.hasOwnProperty(n))continue;this[n]=t[n]}},f=function(e,t){for(var n=0,r=e.length;n<r;n++)if(e[n][0]===t)return n;return-1},l=function(e){var t=e._promiseCallbacks;return t||(t=e._promiseCallbacks={}),t},c={mixin:function(e){return e.on=this.on,e.off=this.off,e.trigger=this.trigger,e},on:function(e,t,n){var r=l(this),i,s;e=e.split(/\s+/),n=n||this;while(s=e.shift())i=r[s],i||(i=r[s]=[]),f(i,t)===-1&&i.push([t,n])},off:function(e,t){var n=l(this),r,i,s;e=e.split(/\s+/);while(i=e.shift()){if(!t){n[i]=[];continue}r=n[i],s=f(r,t),s!==-1&&r.splice(s,1)}},trigger:function(e,t){var n=l(this),r,i,s,o,u;if(r=n[e])for(var f=0,c=r.length;f<c;f++)i=r[f],s=i[0],o=i[1],typeof t!="object"&&(t={detail:t}),u=new a(e,t),s.call(o,u)}},h=function(){this.on("promise:resolved",function(e){this.trigger("success",{detail:e.detail})},this),this.on("promise:failed",function(e){this.trigger("error",{detail:e.detail})},this)},p=function(){},d=function(e,t,n,r){var i=typeof n=="function",s,o,u,a;if(i)try{s=n(r.detail),u=!0}catch(f){a=!0,o=f}else s=r.detail,u=!0;s&&typeof s.then=="function"?s.then(function(e){t.resolve(e)},function(e){t.reject(e)}):i&&u?t.resolve(s):a?t.reject(o):t[e](s)};h.prototype={then:function(e,n){var r=new h;return this.isResolved&&t.async(function(){d("resolve",r,e,{detail:this.resolvedValue})},this),this.isRejected&&t.async(function(){d("reject",r,n,{detail:this.rejectedValue})},this),this.on("promise:resolved",function(t){d("resolve",r,e,t)}),this.on("promise:failed",function(e){d("reject",r,n,e)}),r},resolve:function(e){v(this,e),this.resolve=p,this.reject=p},reject:function(e){m(this,e),this.resolve=p,this.reject=p}},c.mixin(h.prototype),e.Promise=h,e.Event=a,e.EventTarget=c,e.all=g,e.configure=y})(window.RSVP={});