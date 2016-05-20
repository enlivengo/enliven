//! moment.js
//! version : 2.11.1
//! authors : Tim Wood, Iskren Chernev, Moment.js contributors
//! license : MIT
//! momentjs.com
!function(a,b){"object"==typeof exports&&"undefined"!=typeof module?module.exports=b():"function"==typeof define&&define.amd?define(b):a.moment=b()}(this,function(){"use strict";function a(){return Uc.apply(null,arguments)}function b(a){Uc=a}function c(a){return"[object Array]"===Object.prototype.toString.call(a)}function d(a){return a instanceof Date||"[object Date]"===Object.prototype.toString.call(a)}function e(a,b){var c,d=[];for(c=0;c<a.length;++c)d.push(b(a[c],c));return d}function f(a,b){return Object.prototype.hasOwnProperty.call(a,b)}function g(a,b){for(var c in b)f(b,c)&&(a[c]=b[c]);return f(b,"toString")&&(a.toString=b.toString),f(b,"valueOf")&&(a.valueOf=b.valueOf),a}function h(a,b,c,d){return Da(a,b,c,d,!0).utc()}function i(){return{empty:!1,unusedTokens:[],unusedInput:[],overflow:-2,charsLeftOver:0,nullInput:!1,invalidMonth:null,invalidFormat:!1,userInvalidated:!1,iso:!1}}function j(a){return null==a._pf&&(a._pf=i()),a._pf}function k(a){if(null==a._isValid){var b=j(a);a._isValid=!(isNaN(a._d.getTime())||!(b.overflow<0)||b.empty||b.invalidMonth||b.invalidWeekday||b.nullInput||b.invalidFormat||b.userInvalidated),a._strict&&(a._isValid=a._isValid&&0===b.charsLeftOver&&0===b.unusedTokens.length&&void 0===b.bigHour)}return a._isValid}function l(a){var b=h(NaN);return null!=a?g(j(b),a):j(b).userInvalidated=!0,b}function m(a){return void 0===a}function n(a,b){var c,d,e;if(m(b._isAMomentObject)||(a._isAMomentObject=b._isAMomentObject),m(b._i)||(a._i=b._i),m(b._f)||(a._f=b._f),m(b._l)||(a._l=b._l),m(b._strict)||(a._strict=b._strict),m(b._tzm)||(a._tzm=b._tzm),m(b._isUTC)||(a._isUTC=b._isUTC),m(b._offset)||(a._offset=b._offset),m(b._pf)||(a._pf=j(b)),m(b._locale)||(a._locale=b._locale),Wc.length>0)for(c in Wc)d=Wc[c],e=b[d],m(e)||(a[d]=e);return a}function o(b){n(this,b),this._d=new Date(null!=b._d?b._d.getTime():NaN),Xc===!1&&(Xc=!0,a.updateOffset(this),Xc=!1)}function p(a){return a instanceof o||null!=a&&null!=a._isAMomentObject}function q(a){return 0>a?Math.ceil(a):Math.floor(a)}function r(a){var b=+a,c=0;return 0!==b&&isFinite(b)&&(c=q(b)),c}function s(a,b,c){var d,e=Math.min(a.length,b.length),f=Math.abs(a.length-b.length),g=0;for(d=0;e>d;d++)(c&&a[d]!==b[d]||!c&&r(a[d])!==r(b[d]))&&g++;return g+f}function t(){}function u(a){return a?a.toLowerCase().replace("_","-"):a}function v(a){for(var b,c,d,e,f=0;f<a.length;){for(e=u(a[f]).split("-"),b=e.length,c=u(a[f+1]),c=c?c.split("-"):null;b>0;){if(d=w(e.slice(0,b).join("-")))return d;if(c&&c.length>=b&&s(e,c,!0)>=b-1)break;b--}f++}return null}function w(a){var b=null;if(!Yc[a]&&"undefined"!=typeof module&&module&&module.exports)try{b=Vc._abbr,require("./locale/"+a),x(b)}catch(c){}return Yc[a]}function x(a,b){var c;return a&&(c=m(b)?z(a):y(a,b),c&&(Vc=c)),Vc._abbr}function y(a,b){return null!==b?(b.abbr=a,Yc[a]=Yc[a]||new t,Yc[a].set(b),x(a),Yc[a]):(delete Yc[a],null)}function z(a){var b;if(a&&a._locale&&a._locale._abbr&&(a=a._locale._abbr),!a)return Vc;if(!c(a)){if(b=w(a))return b;a=[a]}return v(a)}function A(a,b){var c=a.toLowerCase();Zc[c]=Zc[c+"s"]=Zc[b]=a}function B(a){return"string"==typeof a?Zc[a]||Zc[a.toLowerCase()]:void 0}function C(a){var b,c,d={};for(c in a)f(a,c)&&(b=B(c),b&&(d[b]=a[c]));return d}function D(a){return a instanceof Function||"[object Function]"===Object.prototype.toString.call(a)}function E(b,c){return function(d){return null!=d?(G(this,b,d),a.updateOffset(this,c),this):F(this,b)}}function F(a,b){return a.isValid()?a._d["get"+(a._isUTC?"UTC":"")+b]():NaN}function G(a,b,c){a.isValid()&&a._d["set"+(a._isUTC?"UTC":"")+b](c)}function H(a,b){var c;if("object"==typeof a)for(c in a)this.set(c,a[c]);else if(a=B(a),D(this[a]))return this[a](b);return this}function I(a,b,c){var d=""+Math.abs(a),e=b-d.length,f=a>=0;return(f?c?"+":"":"-")+Math.pow(10,Math.max(0,e)).toString().substr(1)+d}function J(a,b,c,d){var e=d;"string"==typeof d&&(e=function(){return this[d]()}),a&&(bd[a]=e),b&&(bd[b[0]]=function(){return I(e.apply(this,arguments),b[1],b[2])}),c&&(bd[c]=function(){return this.localeData().ordinal(e.apply(this,arguments),a)})}function K(a){return a.match(/\[[\s\S]/)?a.replace(/^\[|\]$/g,""):a.replace(/\\/g,"")}function L(a){var b,c,d=a.match($c);for(b=0,c=d.length;c>b;b++)bd[d[b]]?d[b]=bd[d[b]]:d[b]=K(d[b]);return function(e){var f="";for(b=0;c>b;b++)f+=d[b]instanceof Function?d[b].call(e,a):d[b];return f}}function M(a,b){return a.isValid()?(b=N(b,a.localeData()),ad[b]=ad[b]||L(b),ad[b](a)):a.localeData().invalidDate()}function N(a,b){function c(a){return b.longDateFormat(a)||a}var d=5;for(_c.lastIndex=0;d>=0&&_c.test(a);)a=a.replace(_c,c),_c.lastIndex=0,d-=1;return a}function O(a,b,c){td[a]=D(b)?b:function(a,d){return a&&c?c:b}}function P(a,b){return f(td,a)?td[a](b._strict,b._locale):new RegExp(Q(a))}function Q(a){return R(a.replace("\\","").replace(/\\(\[)|\\(\])|\[([^\]\[]*)\]|\\(.)/g,function(a,b,c,d,e){return b||c||d||e}))}function R(a){return a.replace(/[-\/\\^$*+?.()|[\]{}]/g,"\\$&")}function S(a,b){var c,d=b;for("string"==typeof a&&(a=[a]),"number"==typeof b&&(d=function(a,c){c[b]=r(a)}),c=0;c<a.length;c++)ud[a[c]]=d}function T(a,b){S(a,function(a,c,d,e){d._w=d._w||{},b(a,d._w,d,e)})}function U(a,b,c){null!=b&&f(ud,a)&&ud[a](b,c._a,c,a)}function V(a,b){return new Date(Date.UTC(a,b+1,0)).getUTCDate()}function W(a,b){return c(this._months)?this._months[a.month()]:this._months[Ed.test(b)?"format":"standalone"][a.month()]}function X(a,b){return c(this._monthsShort)?this._monthsShort[a.month()]:this._monthsShort[Ed.test(b)?"format":"standalone"][a.month()]}function Y(a,b,c){var d,e,f;for(this._monthsParse||(this._monthsParse=[],this._longMonthsParse=[],this._shortMonthsParse=[]),d=0;12>d;d++){if(e=h([2e3,d]),c&&!this._longMonthsParse[d]&&(this._longMonthsParse[d]=new RegExp("^"+this.months(e,"").replace(".","")+"$","i"),this._shortMonthsParse[d]=new RegExp("^"+this.monthsShort(e,"").replace(".","")+"$","i")),c||this._monthsParse[d]||(f="^"+this.months(e,"")+"|^"+this.monthsShort(e,""),this._monthsParse[d]=new RegExp(f.replace(".",""),"i")),c&&"MMMM"===b&&this._longMonthsParse[d].test(a))return d;if(c&&"MMM"===b&&this._shortMonthsParse[d].test(a))return d;if(!c&&this._monthsParse[d].test(a))return d}}function Z(a,b){var c;return a.isValid()?"string"==typeof b&&(b=a.localeData().monthsParse(b),"number"!=typeof b)?a:(c=Math.min(a.date(),V(a.year(),b)),a._d["set"+(a._isUTC?"UTC":"")+"Month"](b,c),a):a}function $(b){return null!=b?(Z(this,b),a.updateOffset(this,!0),this):F(this,"Month")}function _(){return V(this.year(),this.month())}function aa(a){return this._monthsParseExact?(f(this,"_monthsRegex")||ca.call(this),a?this._monthsShortStrictRegex:this._monthsShortRegex):this._monthsShortStrictRegex&&a?this._monthsShortStrictRegex:this._monthsShortRegex}function ba(a){return this._monthsParseExact?(f(this,"_monthsRegex")||ca.call(this),a?this._monthsStrictRegex:this._monthsRegex):this._monthsStrictRegex&&a?this._monthsStrictRegex:this._monthsRegex}function ca(){function a(a,b){return b.length-a.length}var b,c,d=[],e=[],f=[];for(b=0;12>b;b++)c=h([2e3,b]),d.push(this.monthsShort(c,"")),e.push(this.months(c,"")),f.push(this.months(c,"")),f.push(this.monthsShort(c,""));for(d.sort(a),e.sort(a),f.sort(a),b=0;12>b;b++)d[b]=R(d[b]),e[b]=R(e[b]),f[b]=R(f[b]);this._monthsRegex=new RegExp("^("+f.join("|")+")","i"),this._monthsShortRegex=this._monthsRegex,this._monthsStrictRegex=new RegExp("^("+e.join("|")+")$","i"),this._monthsShortStrictRegex=new RegExp("^("+d.join("|")+")$","i")}function da(a){var b,c=a._a;return c&&-2===j(a).overflow&&(b=c[wd]<0||c[wd]>11?wd:c[xd]<1||c[xd]>V(c[vd],c[wd])?xd:c[yd]<0||c[yd]>24||24===c[yd]&&(0!==c[zd]||0!==c[Ad]||0!==c[Bd])?yd:c[zd]<0||c[zd]>59?zd:c[Ad]<0||c[Ad]>59?Ad:c[Bd]<0||c[Bd]>999?Bd:-1,j(a)._overflowDayOfYear&&(vd>b||b>xd)&&(b=xd),j(a)._overflowWeeks&&-1===b&&(b=Cd),j(a)._overflowWeekday&&-1===b&&(b=Dd),j(a).overflow=b),a}function ea(b){a.suppressDeprecationWarnings===!1&&"undefined"!=typeof console&&console.warn&&console.warn("Deprecation warning: "+b)}function fa(a,b){var c=!0;return g(function(){return c&&(ea(a+"\nArguments: "+Array.prototype.slice.call(arguments).join(", ")+"\n"+(new Error).stack),c=!1),b.apply(this,arguments)},b)}function ga(a,b){Jd[a]||(ea(b),Jd[a]=!0)}function ha(a){var b,c,d,e,f,g,h=a._i,i=Kd.exec(h)||Ld.exec(h);if(i){for(j(a).iso=!0,b=0,c=Nd.length;c>b;b++)if(Nd[b][1].exec(i[1])){e=Nd[b][0],d=Nd[b][2]!==!1;break}if(null==e)return void(a._isValid=!1);if(i[3]){for(b=0,c=Od.length;c>b;b++)if(Od[b][1].exec(i[3])){f=(i[2]||" ")+Od[b][0];break}if(null==f)return void(a._isValid=!1)}if(!d&&null!=f)return void(a._isValid=!1);if(i[4]){if(!Md.exec(i[4]))return void(a._isValid=!1);g="Z"}a._f=e+(f||"")+(g||""),wa(a)}else a._isValid=!1}function ia(b){var c=Pd.exec(b._i);return null!==c?void(b._d=new Date(+c[1])):(ha(b),void(b._isValid===!1&&(delete b._isValid,a.createFromInputFallback(b))))}function ja(a,b,c,d,e,f,g){var h=new Date(a,b,c,d,e,f,g);return 100>a&&a>=0&&isFinite(h.getFullYear())&&h.setFullYear(a),h}function ka(a){var b=new Date(Date.UTC.apply(null,arguments));return 100>a&&a>=0&&isFinite(b.getUTCFullYear())&&b.setUTCFullYear(a),b}function la(a){return ma(a)?366:365}function ma(a){return a%4===0&&a%100!==0||a%400===0}function na(){return ma(this.year())}function oa(a,b,c){var d=7+b-c,e=(7+ka(a,0,d).getUTCDay()-b)%7;return-e+d-1}function pa(a,b,c,d,e){var f,g,h=(7+c-d)%7,i=oa(a,d,e),j=1+7*(b-1)+h+i;return 0>=j?(f=a-1,g=la(f)+j):j>la(a)?(f=a+1,g=j-la(a)):(f=a,g=j),{year:f,dayOfYear:g}}function qa(a,b,c){var d,e,f=oa(a.year(),b,c),g=Math.floor((a.dayOfYear()-f-1)/7)+1;return 1>g?(e=a.year()-1,d=g+ra(e,b,c)):g>ra(a.year(),b,c)?(d=g-ra(a.year(),b,c),e=a.year()+1):(e=a.year(),d=g),{week:d,year:e}}function ra(a,b,c){var d=oa(a,b,c),e=oa(a+1,b,c);return(la(a)-d+e)/7}function sa(a,b,c){return null!=a?a:null!=b?b:c}function ta(b){var c=new Date(a.now());return b._useUTC?[c.getUTCFullYear(),c.getUTCMonth(),c.getUTCDate()]:[c.getFullYear(),c.getMonth(),c.getDate()]}function ua(a){var b,c,d,e,f=[];if(!a._d){for(d=ta(a),a._w&&null==a._a[xd]&&null==a._a[wd]&&va(a),a._dayOfYear&&(e=sa(a._a[vd],d[vd]),a._dayOfYear>la(e)&&(j(a)._overflowDayOfYear=!0),c=ka(e,0,a._dayOfYear),a._a[wd]=c.getUTCMonth(),a._a[xd]=c.getUTCDate()),b=0;3>b&&null==a._a[b];++b)a._a[b]=f[b]=d[b];for(;7>b;b++)a._a[b]=f[b]=null==a._a[b]?2===b?1:0:a._a[b];24===a._a[yd]&&0===a._a[zd]&&0===a._a[Ad]&&0===a._a[Bd]&&(a._nextDay=!0,a._a[yd]=0),a._d=(a._useUTC?ka:ja).apply(null,f),null!=a._tzm&&a._d.setUTCMinutes(a._d.getUTCMinutes()-a._tzm),a._nextDay&&(a._a[yd]=24)}}function va(a){var b,c,d,e,f,g,h,i;b=a._w,null!=b.GG||null!=b.W||null!=b.E?(f=1,g=4,c=sa(b.GG,a._a[vd],qa(Ea(),1,4).year),d=sa(b.W,1),e=sa(b.E,1),(1>e||e>7)&&(i=!0)):(f=a._locale._week.dow,g=a._locale._week.doy,c=sa(b.gg,a._a[vd],qa(Ea(),f,g).year),d=sa(b.w,1),null!=b.d?(e=b.d,(0>e||e>6)&&(i=!0)):null!=b.e?(e=b.e+f,(b.e<0||b.e>6)&&(i=!0)):e=f),1>d||d>ra(c,f,g)?j(a)._overflowWeeks=!0:null!=i?j(a)._overflowWeekday=!0:(h=pa(c,d,e,f,g),a._a[vd]=h.year,a._dayOfYear=h.dayOfYear)}function wa(b){if(b._f===a.ISO_8601)return void ha(b);b._a=[],j(b).empty=!0;var c,d,e,f,g,h=""+b._i,i=h.length,k=0;for(e=N(b._f,b._locale).match($c)||[],c=0;c<e.length;c++)f=e[c],d=(h.match(P(f,b))||[])[0],d&&(g=h.substr(0,h.indexOf(d)),g.length>0&&j(b).unusedInput.push(g),h=h.slice(h.indexOf(d)+d.length),k+=d.length),bd[f]?(d?j(b).empty=!1:j(b).unusedTokens.push(f),U(f,d,b)):b._strict&&!d&&j(b).unusedTokens.push(f);j(b).charsLeftOver=i-k,h.length>0&&j(b).unusedInput.push(h),j(b).bigHour===!0&&b._a[yd]<=12&&b._a[yd]>0&&(j(b).bigHour=void 0),b._a[yd]=xa(b._locale,b._a[yd],b._meridiem),ua(b),da(b)}function xa(a,b,c){var d;return null==c?b:null!=a.meridiemHour?a.meridiemHour(b,c):null!=a.isPM?(d=a.isPM(c),d&&12>b&&(b+=12),d||12!==b||(b=0),b):b}function ya(a){var b,c,d,e,f;if(0===a._f.length)return j(a).invalidFormat=!0,void(a._d=new Date(NaN));for(e=0;e<a._f.length;e++)f=0,b=n({},a),null!=a._useUTC&&(b._useUTC=a._useUTC),b._f=a._f[e],wa(b),k(b)&&(f+=j(b).charsLeftOver,f+=10*j(b).unusedTokens.length,j(b).score=f,(null==d||d>f)&&(d=f,c=b));g(a,c||b)}function za(a){if(!a._d){var b=C(a._i);a._a=e([b.year,b.month,b.day||b.date,b.hour,b.minute,b.second,b.millisecond],function(a){return a&&parseInt(a,10)}),ua(a)}}function Aa(a){var b=new o(da(Ba(a)));return b._nextDay&&(b.add(1,"d"),b._nextDay=void 0),b}function Ba(a){var b=a._i,e=a._f;return a._locale=a._locale||z(a._l),null===b||void 0===e&&""===b?l({nullInput:!0}):("string"==typeof b&&(a._i=b=a._locale.preparse(b)),p(b)?new o(da(b)):(c(e)?ya(a):e?wa(a):d(b)?a._d=b:Ca(a),k(a)||(a._d=null),a))}function Ca(b){var f=b._i;void 0===f?b._d=new Date(a.now()):d(f)?b._d=new Date(+f):"string"==typeof f?ia(b):c(f)?(b._a=e(f.slice(0),function(a){return parseInt(a,10)}),ua(b)):"object"==typeof f?za(b):"number"==typeof f?b._d=new Date(f):a.createFromInputFallback(b)}function Da(a,b,c,d,e){var f={};return"boolean"==typeof c&&(d=c,c=void 0),f._isAMomentObject=!0,f._useUTC=f._isUTC=e,f._l=c,f._i=a,f._f=b,f._strict=d,Aa(f)}function Ea(a,b,c,d){return Da(a,b,c,d,!1)}function Fa(a,b){var d,e;if(1===b.length&&c(b[0])&&(b=b[0]),!b.length)return Ea();for(d=b[0],e=1;e<b.length;++e)(!b[e].isValid()||b[e][a](d))&&(d=b[e]);return d}function Ga(){var a=[].slice.call(arguments,0);return Fa("isBefore",a)}function Ha(){var a=[].slice.call(arguments,0);return Fa("isAfter",a)}function Ia(a){var b=C(a),c=b.year||0,d=b.quarter||0,e=b.month||0,f=b.week||0,g=b.day||0,h=b.hour||0,i=b.minute||0,j=b.second||0,k=b.millisecond||0;this._milliseconds=+k+1e3*j+6e4*i+36e5*h,this._days=+g+7*f,this._months=+e+3*d+12*c,this._data={},this._locale=z(),this._bubble()}function Ja(a){return a instanceof Ia}function Ka(a,b){J(a,0,0,function(){var a=this.utcOffset(),c="+";return 0>a&&(a=-a,c="-"),c+I(~~(a/60),2)+b+I(~~a%60,2)})}function La(a,b){var c=(b||"").match(a)||[],d=c[c.length-1]||[],e=(d+"").match(Ud)||["-",0,0],f=+(60*e[1])+r(e[2]);return"+"===e[0]?f:-f}function Ma(b,c){var e,f;return c._isUTC?(e=c.clone(),f=(p(b)||d(b)?+b:+Ea(b))-+e,e._d.setTime(+e._d+f),a.updateOffset(e,!1),e):Ea(b).local()}function Na(a){return 15*-Math.round(a._d.getTimezoneOffset()/15)}function Oa(b,c){var d,e=this._offset||0;return this.isValid()?null!=b?("string"==typeof b?b=La(qd,b):Math.abs(b)<16&&(b=60*b),!this._isUTC&&c&&(d=Na(this)),this._offset=b,this._isUTC=!0,null!=d&&this.add(d,"m"),e!==b&&(!c||this._changeInProgress?cb(this,Za(b-e,"m"),1,!1):this._changeInProgress||(this._changeInProgress=!0,a.updateOffset(this,!0),this._changeInProgress=null)),this):this._isUTC?e:Na(this):null!=b?this:NaN}function Pa(a,b){return null!=a?("string"!=typeof a&&(a=-a),this.utcOffset(a,b),this):-this.utcOffset()}function Qa(a){return this.utcOffset(0,a)}function Ra(a){return this._isUTC&&(this.utcOffset(0,a),this._isUTC=!1,a&&this.subtract(Na(this),"m")),this}function Sa(){return this._tzm?this.utcOffset(this._tzm):"string"==typeof this._i&&this.utcOffset(La(pd,this._i)),this}function Ta(a){return this.isValid()?(a=a?Ea(a).utcOffset():0,(this.utcOffset()-a)%60===0):!1}function Ua(){return this.utcOffset()>this.clone().month(0).utcOffset()||this.utcOffset()>this.clone().month(5).utcOffset()}function Va(){if(!m(this._isDSTShifted))return this._isDSTShifted;var a={};if(n(a,this),a=Ba(a),a._a){var b=a._isUTC?h(a._a):Ea(a._a);this._isDSTShifted=this.isValid()&&s(a._a,b.toArray())>0}else this._isDSTShifted=!1;return this._isDSTShifted}function Wa(){return this.isValid()?!this._isUTC:!1}function Xa(){return this.isValid()?this._isUTC:!1}function Ya(){return this.isValid()?this._isUTC&&0===this._offset:!1}function Za(a,b){var c,d,e,g=a,h=null;return Ja(a)?g={ms:a._milliseconds,d:a._days,M:a._months}:"number"==typeof a?(g={},b?g[b]=a:g.milliseconds=a):(h=Vd.exec(a))?(c="-"===h[1]?-1:1,g={y:0,d:r(h[xd])*c,h:r(h[yd])*c,m:r(h[zd])*c,s:r(h[Ad])*c,ms:r(h[Bd])*c}):(h=Wd.exec(a))?(c="-"===h[1]?-1:1,g={y:$a(h[2],c),M:$a(h[3],c),d:$a(h[4],c),h:$a(h[5],c),m:$a(h[6],c),s:$a(h[7],c),w:$a(h[8],c)}):null==g?g={}:"object"==typeof g&&("from"in g||"to"in g)&&(e=ab(Ea(g.from),Ea(g.to)),g={},g.ms=e.milliseconds,g.M=e.months),d=new Ia(g),Ja(a)&&f(a,"_locale")&&(d._locale=a._locale),d}function $a(a,b){var c=a&&parseFloat(a.replace(",","."));return(isNaN(c)?0:c)*b}function _a(a,b){var c={milliseconds:0,months:0};return c.months=b.month()-a.month()+12*(b.year()-a.year()),a.clone().add(c.months,"M").isAfter(b)&&--c.months,c.milliseconds=+b-+a.clone().add(c.months,"M"),c}function ab(a,b){var c;return a.isValid()&&b.isValid()?(b=Ma(b,a),a.isBefore(b)?c=_a(a,b):(c=_a(b,a),c.milliseconds=-c.milliseconds,c.months=-c.months),c):{milliseconds:0,months:0}}function bb(a,b){return function(c,d){var e,f;return null===d||isNaN(+d)||(ga(b,"moment()."+b+"(period, number) is deprecated. Please use moment()."+b+"(number, period)."),f=c,c=d,d=f),c="string"==typeof c?+c:c,e=Za(c,d),cb(this,e,a),this}}function cb(b,c,d,e){var f=c._milliseconds,g=c._days,h=c._months;b.isValid()&&(e=null==e?!0:e,f&&b._d.setTime(+b._d+f*d),g&&G(b,"Date",F(b,"Date")+g*d),h&&Z(b,F(b,"Month")+h*d),e&&a.updateOffset(b,g||h))}function db(a,b){var c=a||Ea(),d=Ma(c,this).startOf("day"),e=this.diff(d,"days",!0),f=-6>e?"sameElse":-1>e?"lastWeek":0>e?"lastDay":1>e?"sameDay":2>e?"nextDay":7>e?"nextWeek":"sameElse",g=b&&(D(b[f])?b[f]():b[f]);return this.format(g||this.localeData().calendar(f,this,Ea(c)))}function eb(){return new o(this)}function fb(a,b){var c=p(a)?a:Ea(a);return this.isValid()&&c.isValid()?(b=B(m(b)?"millisecond":b),"millisecond"===b?+this>+c:+c<+this.clone().startOf(b)):!1}function gb(a,b){var c=p(a)?a:Ea(a);return this.isValid()&&c.isValid()?(b=B(m(b)?"millisecond":b),"millisecond"===b?+c>+this:+this.clone().endOf(b)<+c):!1}function hb(a,b,c){return this.isAfter(a,c)&&this.isBefore(b,c)}function ib(a,b){var c,d=p(a)?a:Ea(a);return this.isValid()&&d.isValid()?(b=B(b||"millisecond"),"millisecond"===b?+this===+d:(c=+d,+this.clone().startOf(b)<=c&&c<=+this.clone().endOf(b))):!1}function jb(a,b){return this.isSame(a,b)||this.isAfter(a,b)}function kb(a,b){return this.isSame(a,b)||this.isBefore(a,b)}function lb(a,b,c){var d,e,f,g;return this.isValid()?(d=Ma(a,this),d.isValid()?(e=6e4*(d.utcOffset()-this.utcOffset()),b=B(b),"year"===b||"month"===b||"quarter"===b?(g=mb(this,d),"quarter"===b?g/=3:"year"===b&&(g/=12)):(f=this-d,g="second"===b?f/1e3:"minute"===b?f/6e4:"hour"===b?f/36e5:"day"===b?(f-e)/864e5:"week"===b?(f-e)/6048e5:f),c?g:q(g)):NaN):NaN}function mb(a,b){var c,d,e=12*(b.year()-a.year())+(b.month()-a.month()),f=a.clone().add(e,"months");return 0>b-f?(c=a.clone().add(e-1,"months"),d=(b-f)/(f-c)):(c=a.clone().add(e+1,"months"),d=(b-f)/(c-f)),-(e+d)}function nb(){return this.clone().locale("en").format("ddd MMM DD YYYY HH:mm:ss [GMT]ZZ")}function ob(){var a=this.clone().utc();return 0<a.year()&&a.year()<=9999?D(Date.prototype.toISOString)?this.toDate().toISOString():M(a,"YYYY-MM-DD[T]HH:mm:ss.SSS[Z]"):M(a,"YYYYYY-MM-DD[T]HH:mm:ss.SSS[Z]")}function pb(b){var c=M(this,b||a.defaultFormat);return this.localeData().postformat(c)}function qb(a,b){return this.isValid()&&(p(a)&&a.isValid()||Ea(a).isValid())?Za({to:this,from:a}).locale(this.locale()).humanize(!b):this.localeData().invalidDate()}function rb(a){return this.from(Ea(),a)}function sb(a,b){return this.isValid()&&(p(a)&&a.isValid()||Ea(a).isValid())?Za({from:this,to:a}).locale(this.locale()).humanize(!b):this.localeData().invalidDate()}function tb(a){return this.to(Ea(),a)}function ub(a){var b;return void 0===a?this._locale._abbr:(b=z(a),null!=b&&(this._locale=b),this)}function vb(){return this._locale}function wb(a){switch(a=B(a)){case"year":this.month(0);case"quarter":case"month":this.date(1);case"week":case"isoWeek":case"day":this.hours(0);case"hour":this.minutes(0);case"minute":this.seconds(0);case"second":this.milliseconds(0)}return"week"===a&&this.weekday(0),"isoWeek"===a&&this.isoWeekday(1),"quarter"===a&&this.month(3*Math.floor(this.month()/3)),this}function xb(a){return a=B(a),void 0===a||"millisecond"===a?this:this.startOf(a).add(1,"isoWeek"===a?"week":a).subtract(1,"ms")}function yb(){return+this._d-6e4*(this._offset||0)}function zb(){return Math.floor(+this/1e3)}function Ab(){return this._offset?new Date(+this):this._d}function Bb(){var a=this;return[a.year(),a.month(),a.date(),a.hour(),a.minute(),a.second(),a.millisecond()]}function Cb(){var a=this;return{years:a.year(),months:a.month(),date:a.date(),hours:a.hours(),minutes:a.minutes(),seconds:a.seconds(),milliseconds:a.milliseconds()}}function Db(){return this.isValid()?this.toISOString():"null"}function Eb(){return k(this)}function Fb(){return g({},j(this))}function Gb(){return j(this).overflow}function Hb(){return{input:this._i,format:this._f,locale:this._locale,isUTC:this._isUTC,strict:this._strict}}function Ib(a,b){J(0,[a,a.length],0,b)}function Jb(a){return Nb.call(this,a,this.week(),this.weekday(),this.localeData()._week.dow,this.localeData()._week.doy)}function Kb(a){return Nb.call(this,a,this.isoWeek(),this.isoWeekday(),1,4)}function Lb(){return ra(this.year(),1,4)}function Mb(){var a=this.localeData()._week;return ra(this.year(),a.dow,a.doy)}function Nb(a,b,c,d,e){var f;return null==a?qa(this,d,e).year:(f=ra(a,d,e),b>f&&(b=f),Ob.call(this,a,b,c,d,e))}function Ob(a,b,c,d,e){var f=pa(a,b,c,d,e),g=ka(f.year,0,f.dayOfYear);return this.year(g.getUTCFullYear()),this.month(g.getUTCMonth()),this.date(g.getUTCDate()),this}function Pb(a){return null==a?Math.ceil((this.month()+1)/3):this.month(3*(a-1)+this.month()%3)}function Qb(a){return qa(a,this._week.dow,this._week.doy).week}function Rb(){return this._week.dow}function Sb(){return this._week.doy}function Tb(a){var b=this.localeData().week(this);return null==a?b:this.add(7*(a-b),"d")}function Ub(a){var b=qa(this,1,4).week;return null==a?b:this.add(7*(a-b),"d")}function Vb(a,b){return"string"!=typeof a?a:isNaN(a)?(a=b.weekdaysParse(a),"number"==typeof a?a:null):parseInt(a,10)}function Wb(a,b){return c(this._weekdays)?this._weekdays[a.day()]:this._weekdays[this._weekdays.isFormat.test(b)?"format":"standalone"][a.day()]}function Xb(a){return this._weekdaysShort[a.day()]}function Yb(a){return this._weekdaysMin[a.day()]}function Zb(a,b,c){var d,e,f;for(this._weekdaysParse||(this._weekdaysParse=[],this._minWeekdaysParse=[],this._shortWeekdaysParse=[],this._fullWeekdaysParse=[]),d=0;7>d;d++){if(e=Ea([2e3,1]).day(d),c&&!this._fullWeekdaysParse[d]&&(this._fullWeekdaysParse[d]=new RegExp("^"+this.weekdays(e,"").replace(".",".?")+"$","i"),this._shortWeekdaysParse[d]=new RegExp("^"+this.weekdaysShort(e,"").replace(".",".?")+"$","i"),this._minWeekdaysParse[d]=new RegExp("^"+this.weekdaysMin(e,"").replace(".",".?")+"$","i")),this._weekdaysParse[d]||(f="^"+this.weekdays(e,"")+"|^"+this.weekdaysShort(e,"")+"|^"+this.weekdaysMin(e,""),this._weekdaysParse[d]=new RegExp(f.replace(".",""),"i")),c&&"dddd"===b&&this._fullWeekdaysParse[d].test(a))return d;if(c&&"ddd"===b&&this._shortWeekdaysParse[d].test(a))return d;if(c&&"dd"===b&&this._minWeekdaysParse[d].test(a))return d;if(!c&&this._weekdaysParse[d].test(a))return d}}function $b(a){if(!this.isValid())return null!=a?this:NaN;var b=this._isUTC?this._d.getUTCDay():this._d.getDay();return null!=a?(a=Vb(a,this.localeData()),this.add(a-b,"d")):b}function _b(a){if(!this.isValid())return null!=a?this:NaN;var b=(this.day()+7-this.localeData()._week.dow)%7;return null==a?b:this.add(a-b,"d")}function ac(a){return this.isValid()?null==a?this.day()||7:this.day(this.day()%7?a:a-7):null!=a?this:NaN}function bc(a){var b=Math.round((this.clone().startOf("day")-this.clone().startOf("year"))/864e5)+1;return null==a?b:this.add(a-b,"d")}function cc(){return this.hours()%12||12}function dc(a,b){J(a,0,0,function(){return this.localeData().meridiem(this.hours(),this.minutes(),b)})}function ec(a,b){return b._meridiemParse}function fc(a){return"p"===(a+"").toLowerCase().charAt(0)}function gc(a,b,c){return a>11?c?"pm":"PM":c?"am":"AM"}function hc(a,b){b[Bd]=r(1e3*("0."+a))}function ic(){return this._isUTC?"UTC":""}function jc(){return this._isUTC?"Coordinated Universal Time":""}function kc(a){return Ea(1e3*a)}function lc(){return Ea.apply(null,arguments).parseZone()}function mc(a,b,c){var d=this._calendar[a];return D(d)?d.call(b,c):d}function nc(a){var b=this._longDateFormat[a],c=this._longDateFormat[a.toUpperCase()];return b||!c?b:(this._longDateFormat[a]=c.replace(/MMMM|MM|DD|dddd/g,function(a){return a.slice(1)}),this._longDateFormat[a])}function oc(){return this._invalidDate}function pc(a){return this._ordinal.replace("%d",a)}function qc(a){return a}function rc(a,b,c,d){var e=this._relativeTime[c];return D(e)?e(a,b,c,d):e.replace(/%d/i,a)}function sc(a,b){var c=this._relativeTime[a>0?"future":"past"];return D(c)?c(b):c.replace(/%s/i,b)}function tc(a){var b,c;for(c in a)b=a[c],D(b)?this[c]=b:this["_"+c]=b;this._ordinalParseLenient=new RegExp(this._ordinalParse.source+"|"+/\d{1,2}/.source)}function uc(a,b,c,d){var e=z(),f=h().set(d,b);return e[c](f,a)}function vc(a,b,c,d,e){if("number"==typeof a&&(b=a,a=void 0),a=a||"",null!=b)return uc(a,b,c,e);var f,g=[];for(f=0;d>f;f++)g[f]=uc(a,f,c,e);return g}function wc(a,b){return vc(a,b,"months",12,"month")}function xc(a,b){return vc(a,b,"monthsShort",12,"month")}function yc(a,b){return vc(a,b,"weekdays",7,"day")}function zc(a,b){return vc(a,b,"weekdaysShort",7,"day")}function Ac(a,b){return vc(a,b,"weekdaysMin",7,"day")}function Bc(){var a=this._data;return this._milliseconds=se(this._milliseconds),this._days=se(this._days),this._months=se(this._months),a.milliseconds=se(a.milliseconds),a.seconds=se(a.seconds),a.minutes=se(a.minutes),a.hours=se(a.hours),a.months=se(a.months),a.years=se(a.years),this}function Cc(a,b,c,d){var e=Za(b,c);return a._milliseconds+=d*e._milliseconds,a._days+=d*e._days,a._months+=d*e._months,a._bubble()}function Dc(a,b){return Cc(this,a,b,1)}function Ec(a,b){return Cc(this,a,b,-1)}function Fc(a){return 0>a?Math.floor(a):Math.ceil(a)}function Gc(){var a,b,c,d,e,f=this._milliseconds,g=this._days,h=this._months,i=this._data;return f>=0&&g>=0&&h>=0||0>=f&&0>=g&&0>=h||(f+=864e5*Fc(Ic(h)+g),g=0,h=0),i.milliseconds=f%1e3,a=q(f/1e3),i.seconds=a%60,b=q(a/60),i.minutes=b%60,c=q(b/60),i.hours=c%24,g+=q(c/24),e=q(Hc(g)),h+=e,g-=Fc(Ic(e)),d=q(h/12),h%=12,i.days=g,i.months=h,i.years=d,this}function Hc(a){return 4800*a/146097}function Ic(a){return 146097*a/4800}function Jc(a){var b,c,d=this._milliseconds;if(a=B(a),"month"===a||"year"===a)return b=this._days+d/864e5,c=this._months+Hc(b),"month"===a?c:c/12;switch(b=this._days+Math.round(Ic(this._months)),a){case"week":return b/7+d/6048e5;case"day":return b+d/864e5;case"hour":return 24*b+d/36e5;case"minute":return 1440*b+d/6e4;case"second":return 86400*b+d/1e3;case"millisecond":return Math.floor(864e5*b)+d;default:throw new Error("Unknown unit "+a)}}function Kc(){return this._milliseconds+864e5*this._days+this._months%12*2592e6+31536e6*r(this._months/12)}function Lc(a){return function(){return this.as(a)}}function Mc(a){return a=B(a),this[a+"s"]()}function Nc(a){return function(){return this._data[a]}}function Oc(){return q(this.days()/7)}function Pc(a,b,c,d,e){return e.relativeTime(b||1,!!c,a,d)}function Qc(a,b,c){var d=Za(a).abs(),e=Ie(d.as("s")),f=Ie(d.as("m")),g=Ie(d.as("h")),h=Ie(d.as("d")),i=Ie(d.as("M")),j=Ie(d.as("y")),k=e<Je.s&&["s",e]||1>=f&&["m"]||f<Je.m&&["mm",f]||1>=g&&["h"]||g<Je.h&&["hh",g]||1>=h&&["d"]||h<Je.d&&["dd",h]||1>=i&&["M"]||i<Je.M&&["MM",i]||1>=j&&["y"]||["yy",j];return k[2]=b,k[3]=+a>0,k[4]=c,Pc.apply(null,k)}function Rc(a,b){return void 0===Je[a]?!1:void 0===b?Je[a]:(Je[a]=b,!0)}function Sc(a){var b=this.localeData(),c=Qc(this,!a,b);return a&&(c=b.pastFuture(+this,c)),b.postformat(c)}function Tc(){var a,b,c,d=Ke(this._milliseconds)/1e3,e=Ke(this._days),f=Ke(this._months);a=q(d/60),b=q(a/60),d%=60,a%=60,c=q(f/12),f%=12;var g=c,h=f,i=e,j=b,k=a,l=d,m=this.asSeconds();return m?(0>m?"-":"")+"P"+(g?g+"Y":"")+(h?h+"M":"")+(i?i+"D":"")+(j||k||l?"T":"")+(j?j+"H":"")+(k?k+"M":"")+(l?l+"S":""):"P0D"}var Uc,Vc,Wc=a.momentProperties=[],Xc=!1,Yc={},Zc={},$c=/(\[[^\[]*\])|(\\)?([Hh]mm(ss)?|Mo|MM?M?M?|Do|DDDo|DD?D?D?|ddd?d?|do?|w[o|w]?|W[o|W]?|Qo?|YYYYYY|YYYYY|YYYY|YY|gg(ggg?)?|GG(GGG?)?|e|E|a|A|hh?|HH?|mm?|ss?|S{1,9}|x|X|zz?|ZZ?|.)/g,_c=/(\[[^\[]*\])|(\\)?(LTS|LT|LL?L?L?|l{1,4})/g,ad={},bd={},cd=/\d/,dd=/\d\d/,ed=/\d{3}/,fd=/\d{4}/,gd=/[+-]?\d{6}/,hd=/\d\d?/,id=/\d\d\d\d?/,jd=/\d\d\d\d\d\d?/,kd=/\d{1,3}/,ld=/\d{1,4}/,md=/[+-]?\d{1,6}/,nd=/\d+/,od=/[+-]?\d+/,pd=/Z|[+-]\d\d:?\d\d/gi,qd=/Z|[+-]\d\d(?::?\d\d)?/gi,rd=/[+-]?\d+(\.\d{1,3})?/,sd=/[0-9]*['a-z\u00A0-\u05FF\u0700-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF]+|[\u0600-\u06FF\/]+(\s*?[\u0600-\u06FF]+){1,2}/i,td={},ud={},vd=0,wd=1,xd=2,yd=3,zd=4,Ad=5,Bd=6,Cd=7,Dd=8;J("M",["MM",2],"Mo",function(){return this.month()+1}),J("MMM",0,0,function(a){return this.localeData().monthsShort(this,a)}),J("MMMM",0,0,function(a){return this.localeData().months(this,a)}),A("month","M"),O("M",hd),O("MM",hd,dd),O("MMM",function(a,b){return b.monthsShortRegex(a)}),O("MMMM",function(a,b){return b.monthsRegex(a)}),S(["M","MM"],function(a,b){b[wd]=r(a)-1}),S(["MMM","MMMM"],function(a,b,c,d){var e=c._locale.monthsParse(a,d,c._strict);null!=e?b[wd]=e:j(c).invalidMonth=a});var Ed=/D[oD]?(\[[^\[\]]*\]|\s+)+MMMM?/,Fd="January_February_March_April_May_June_July_August_September_October_November_December".split("_"),Gd="Jan_Feb_Mar_Apr_May_Jun_Jul_Aug_Sep_Oct_Nov_Dec".split("_"),Hd=sd,Id=sd,Jd={};a.suppressDeprecationWarnings=!1;var Kd=/^\s*((?:[+-]\d{6}|\d{4})-(?:\d\d-\d\d|W\d\d-\d|W\d\d|\d\d\d|\d\d))(?:(T| )(\d\d(?::\d\d(?::\d\d(?:[.,]\d+)?)?)?)([\+\-]\d\d(?::?\d\d)?|\s*Z)?)?/,Ld=/^\s*((?:[+-]\d{6}|\d{4})(?:\d\d\d\d|W\d\d\d|W\d\d|\d\d\d|\d\d))(?:(T| )(\d\d(?:\d\d(?:\d\d(?:[.,]\d+)?)?)?)([\+\-]\d\d(?::?\d\d)?|\s*Z)?)?/,Md=/Z|[+-]\d\d(?::?\d\d)?/,Nd=[["YYYYYY-MM-DD",/[+-]\d{6}-\d\d-\d\d/],["YYYY-MM-DD",/\d{4}-\d\d-\d\d/],["GGGG-[W]WW-E",/\d{4}-W\d\d-\d/],["GGGG-[W]WW",/\d{4}-W\d\d/,!1],["YYYY-DDD",/\d{4}-\d{3}/],["YYYY-MM",/\d{4}-\d\d/,!1],["YYYYYYMMDD",/[+-]\d{10}/],["YYYYMMDD",/\d{8}/],["GGGG[W]WWE",/\d{4}W\d{3}/],["GGGG[W]WW",/\d{4}W\d{2}/,!1],["YYYYDDD",/\d{7}/]],Od=[["HH:mm:ss.SSSS",/\d\d:\d\d:\d\d\.\d+/],["HH:mm:ss,SSSS",/\d\d:\d\d:\d\d,\d+/],["HH:mm:ss",/\d\d:\d\d:\d\d/],["HH:mm",/\d\d:\d\d/],["HHmmss.SSSS",/\d\d\d\d\d\d\.\d+/],["HHmmss,SSSS",/\d\d\d\d\d\d,\d+/],["HHmmss",/\d\d\d\d\d\d/],["HHmm",/\d\d\d\d/],["HH",/\d\d/]],Pd=/^\/?Date\((\-?\d+)/i;a.createFromInputFallback=fa("moment construction falls back to js Date. This is discouraged and will be removed in upcoming major release. Please refer to https://github.com/moment/moment/issues/1407 for more info.",function(a){a._d=new Date(a._i+(a._useUTC?" UTC":""))}),J("Y",0,0,function(){var a=this.year();return 9999>=a?""+a:"+"+a}),J(0,["YY",2],0,function(){return this.year()%100}),J(0,["YYYY",4],0,"year"),J(0,["YYYYY",5],0,"year"),J(0,["YYYYYY",6,!0],0,"year"),A("year","y"),O("Y",od),O("YY",hd,dd),O("YYYY",ld,fd),O("YYYYY",md,gd),O("YYYYYY",md,gd),S(["YYYYY","YYYYYY"],vd),S("YYYY",function(b,c){c[vd]=2===b.length?a.parseTwoDigitYear(b):r(b)}),S("YY",function(b,c){c[vd]=a.parseTwoDigitYear(b)}),S("Y",function(a,b){b[vd]=parseInt(a,10)}),a.parseTwoDigitYear=function(a){return r(a)+(r(a)>68?1900:2e3)};var Qd=E("FullYear",!1);a.ISO_8601=function(){};var Rd=fa("moment().min is deprecated, use moment.min instead. https://github.com/moment/moment/issues/1548",function(){var a=Ea.apply(null,arguments);return this.isValid()&&a.isValid()?this>a?this:a:l()}),Sd=fa("moment().max is deprecated, use moment.max instead. https://github.com/moment/moment/issues/1548",function(){var a=Ea.apply(null,arguments);return this.isValid()&&a.isValid()?a>this?this:a:l()}),Td=function(){return Date.now?Date.now():+new Date};Ka("Z",":"),Ka("ZZ",""),O("Z",qd),O("ZZ",qd),S(["Z","ZZ"],function(a,b,c){c._useUTC=!0,c._tzm=La(qd,a)});var Ud=/([\+\-]|\d\d)/gi;a.updateOffset=function(){};var Vd=/(\-)?(?:(\d*)[. ])?(\d+)\:(\d+)(?:\:(\d+)\.?(\d{3})?)?/,Wd=/^(-)?P(?:(?:([0-9,.]*)Y)?(?:([0-9,.]*)M)?(?:([0-9,.]*)D)?(?:T(?:([0-9,.]*)H)?(?:([0-9,.]*)M)?(?:([0-9,.]*)S)?)?|([0-9,.]*)W)$/;
Za.fn=Ia.prototype;var Xd=bb(1,"add"),Yd=bb(-1,"subtract");a.defaultFormat="YYYY-MM-DDTHH:mm:ssZ";var Zd=fa("moment().lang() is deprecated. Instead, use moment().localeData() to get the language configuration. Use moment().locale() to change languages.",function(a){return void 0===a?this.localeData():this.locale(a)});J(0,["gg",2],0,function(){return this.weekYear()%100}),J(0,["GG",2],0,function(){return this.isoWeekYear()%100}),Ib("gggg","weekYear"),Ib("ggggg","weekYear"),Ib("GGGG","isoWeekYear"),Ib("GGGGG","isoWeekYear"),A("weekYear","gg"),A("isoWeekYear","GG"),O("G",od),O("g",od),O("GG",hd,dd),O("gg",hd,dd),O("GGGG",ld,fd),O("gggg",ld,fd),O("GGGGG",md,gd),O("ggggg",md,gd),T(["gggg","ggggg","GGGG","GGGGG"],function(a,b,c,d){b[d.substr(0,2)]=r(a)}),T(["gg","GG"],function(b,c,d,e){c[e]=a.parseTwoDigitYear(b)}),J("Q",0,"Qo","quarter"),A("quarter","Q"),O("Q",cd),S("Q",function(a,b){b[wd]=3*(r(a)-1)}),J("w",["ww",2],"wo","week"),J("W",["WW",2],"Wo","isoWeek"),A("week","w"),A("isoWeek","W"),O("w",hd),O("ww",hd,dd),O("W",hd),O("WW",hd,dd),T(["w","ww","W","WW"],function(a,b,c,d){b[d.substr(0,1)]=r(a)});var $d={dow:0,doy:6};J("D",["DD",2],"Do","date"),A("date","D"),O("D",hd),O("DD",hd,dd),O("Do",function(a,b){return a?b._ordinalParse:b._ordinalParseLenient}),S(["D","DD"],xd),S("Do",function(a,b){b[xd]=r(a.match(hd)[0],10)});var _d=E("Date",!0);J("d",0,"do","day"),J("dd",0,0,function(a){return this.localeData().weekdaysMin(this,a)}),J("ddd",0,0,function(a){return this.localeData().weekdaysShort(this,a)}),J("dddd",0,0,function(a){return this.localeData().weekdays(this,a)}),J("e",0,0,"weekday"),J("E",0,0,"isoWeekday"),A("day","d"),A("weekday","e"),A("isoWeekday","E"),O("d",hd),O("e",hd),O("E",hd),O("dd",sd),O("ddd",sd),O("dddd",sd),T(["dd","ddd","dddd"],function(a,b,c,d){var e=c._locale.weekdaysParse(a,d,c._strict);null!=e?b.d=e:j(c).invalidWeekday=a}),T(["d","e","E"],function(a,b,c,d){b[d]=r(a)});var ae="Sunday_Monday_Tuesday_Wednesday_Thursday_Friday_Saturday".split("_"),be="Sun_Mon_Tue_Wed_Thu_Fri_Sat".split("_"),ce="Su_Mo_Tu_We_Th_Fr_Sa".split("_");J("DDD",["DDDD",3],"DDDo","dayOfYear"),A("dayOfYear","DDD"),O("DDD",kd),O("DDDD",ed),S(["DDD","DDDD"],function(a,b,c){c._dayOfYear=r(a)}),J("H",["HH",2],0,"hour"),J("h",["hh",2],0,cc),J("hmm",0,0,function(){return""+cc.apply(this)+I(this.minutes(),2)}),J("hmmss",0,0,function(){return""+cc.apply(this)+I(this.minutes(),2)+I(this.seconds(),2)}),J("Hmm",0,0,function(){return""+this.hours()+I(this.minutes(),2)}),J("Hmmss",0,0,function(){return""+this.hours()+I(this.minutes(),2)+I(this.seconds(),2)}),dc("a",!0),dc("A",!1),A("hour","h"),O("a",ec),O("A",ec),O("H",hd),O("h",hd),O("HH",hd,dd),O("hh",hd,dd),O("hmm",id),O("hmmss",jd),O("Hmm",id),O("Hmmss",jd),S(["H","HH"],yd),S(["a","A"],function(a,b,c){c._isPm=c._locale.isPM(a),c._meridiem=a}),S(["h","hh"],function(a,b,c){b[yd]=r(a),j(c).bigHour=!0}),S("hmm",function(a,b,c){var d=a.length-2;b[yd]=r(a.substr(0,d)),b[zd]=r(a.substr(d)),j(c).bigHour=!0}),S("hmmss",function(a,b,c){var d=a.length-4,e=a.length-2;b[yd]=r(a.substr(0,d)),b[zd]=r(a.substr(d,2)),b[Ad]=r(a.substr(e)),j(c).bigHour=!0}),S("Hmm",function(a,b,c){var d=a.length-2;b[yd]=r(a.substr(0,d)),b[zd]=r(a.substr(d))}),S("Hmmss",function(a,b,c){var d=a.length-4,e=a.length-2;b[yd]=r(a.substr(0,d)),b[zd]=r(a.substr(d,2)),b[Ad]=r(a.substr(e))});var de=/[ap]\.?m?\.?/i,ee=E("Hours",!0);J("m",["mm",2],0,"minute"),A("minute","m"),O("m",hd),O("mm",hd,dd),S(["m","mm"],zd);var fe=E("Minutes",!1);J("s",["ss",2],0,"second"),A("second","s"),O("s",hd),O("ss",hd,dd),S(["s","ss"],Ad);var ge=E("Seconds",!1);J("S",0,0,function(){return~~(this.millisecond()/100)}),J(0,["SS",2],0,function(){return~~(this.millisecond()/10)}),J(0,["SSS",3],0,"millisecond"),J(0,["SSSS",4],0,function(){return 10*this.millisecond()}),J(0,["SSSSS",5],0,function(){return 100*this.millisecond()}),J(0,["SSSSSS",6],0,function(){return 1e3*this.millisecond()}),J(0,["SSSSSSS",7],0,function(){return 1e4*this.millisecond()}),J(0,["SSSSSSSS",8],0,function(){return 1e5*this.millisecond()}),J(0,["SSSSSSSSS",9],0,function(){return 1e6*this.millisecond()}),A("millisecond","ms"),O("S",kd,cd),O("SS",kd,dd),O("SSS",kd,ed);var he;for(he="SSSS";he.length<=9;he+="S")O(he,nd);for(he="S";he.length<=9;he+="S")S(he,hc);var ie=E("Milliseconds",!1);J("z",0,0,"zoneAbbr"),J("zz",0,0,"zoneName");var je=o.prototype;je.add=Xd,je.calendar=db,je.clone=eb,je.diff=lb,je.endOf=xb,je.format=pb,je.from=qb,je.fromNow=rb,je.to=sb,je.toNow=tb,je.get=H,je.invalidAt=Gb,je.isAfter=fb,je.isBefore=gb,je.isBetween=hb,je.isSame=ib,je.isSameOrAfter=jb,je.isSameOrBefore=kb,je.isValid=Eb,je.lang=Zd,je.locale=ub,je.localeData=vb,je.max=Sd,je.min=Rd,je.parsingFlags=Fb,je.set=H,je.startOf=wb,je.subtract=Yd,je.toArray=Bb,je.toObject=Cb,je.toDate=Ab,je.toISOString=ob,je.toJSON=Db,je.toString=nb,je.unix=zb,je.valueOf=yb,je.creationData=Hb,je.year=Qd,je.isLeapYear=na,je.weekYear=Jb,je.isoWeekYear=Kb,je.quarter=je.quarters=Pb,je.month=$,je.daysInMonth=_,je.week=je.weeks=Tb,je.isoWeek=je.isoWeeks=Ub,je.weeksInYear=Mb,je.isoWeeksInYear=Lb,je.date=_d,je.day=je.days=$b,je.weekday=_b,je.isoWeekday=ac,je.dayOfYear=bc,je.hour=je.hours=ee,je.minute=je.minutes=fe,je.second=je.seconds=ge,je.millisecond=je.milliseconds=ie,je.utcOffset=Oa,je.utc=Qa,je.local=Ra,je.parseZone=Sa,je.hasAlignedHourOffset=Ta,je.isDST=Ua,je.isDSTShifted=Va,je.isLocal=Wa,je.isUtcOffset=Xa,je.isUtc=Ya,je.isUTC=Ya,je.zoneAbbr=ic,je.zoneName=jc,je.dates=fa("dates accessor is deprecated. Use date instead.",_d),je.months=fa("months accessor is deprecated. Use month instead",$),je.years=fa("years accessor is deprecated. Use year instead",Qd),je.zone=fa("moment().zone is deprecated, use moment().utcOffset instead. https://github.com/moment/moment/issues/1779",Pa);var ke=je,le={sameDay:"[Today at] LT",nextDay:"[Tomorrow at] LT",nextWeek:"dddd [at] LT",lastDay:"[Yesterday at] LT",lastWeek:"[Last] dddd [at] LT",sameElse:"L"},me={LTS:"h:mm:ss A",LT:"h:mm A",L:"MM/DD/YYYY",LL:"MMMM D, YYYY",LLL:"MMMM D, YYYY h:mm A",LLLL:"dddd, MMMM D, YYYY h:mm A"},ne="Invalid date",oe="%d",pe=/\d{1,2}/,qe={future:"in %s",past:"%s ago",s:"a few seconds",m:"a minute",mm:"%d minutes",h:"an hour",hh:"%d hours",d:"a day",dd:"%d days",M:"a month",MM:"%d months",y:"a year",yy:"%d years"},re=t.prototype;re._calendar=le,re.calendar=mc,re._longDateFormat=me,re.longDateFormat=nc,re._invalidDate=ne,re.invalidDate=oc,re._ordinal=oe,re.ordinal=pc,re._ordinalParse=pe,re.preparse=qc,re.postformat=qc,re._relativeTime=qe,re.relativeTime=rc,re.pastFuture=sc,re.set=tc,re.months=W,re._months=Fd,re.monthsShort=X,re._monthsShort=Gd,re.monthsParse=Y,re._monthsRegex=Id,re.monthsRegex=ba,re._monthsShortRegex=Hd,re.monthsShortRegex=aa,re.week=Qb,re._week=$d,re.firstDayOfYear=Sb,re.firstDayOfWeek=Rb,re.weekdays=Wb,re._weekdays=ae,re.weekdaysMin=Yb,re._weekdaysMin=ce,re.weekdaysShort=Xb,re._weekdaysShort=be,re.weekdaysParse=Zb,re.isPM=fc,re._meridiemParse=de,re.meridiem=gc,x("en",{ordinalParse:/\d{1,2}(th|st|nd|rd)/,ordinal:function(a){var b=a%10,c=1===r(a%100/10)?"th":1===b?"st":2===b?"nd":3===b?"rd":"th";return a+c}}),a.lang=fa("moment.lang is deprecated. Use moment.locale instead.",x),a.langData=fa("moment.langData is deprecated. Use moment.localeData instead.",z);var se=Math.abs,te=Lc("ms"),ue=Lc("s"),ve=Lc("m"),we=Lc("h"),xe=Lc("d"),ye=Lc("w"),ze=Lc("M"),Ae=Lc("y"),Be=Nc("milliseconds"),Ce=Nc("seconds"),De=Nc("minutes"),Ee=Nc("hours"),Fe=Nc("days"),Ge=Nc("months"),He=Nc("years"),Ie=Math.round,Je={s:45,m:45,h:22,d:26,M:11},Ke=Math.abs,Le=Ia.prototype;Le.abs=Bc,Le.add=Dc,Le.subtract=Ec,Le.as=Jc,Le.asMilliseconds=te,Le.asSeconds=ue,Le.asMinutes=ve,Le.asHours=we,Le.asDays=xe,Le.asWeeks=ye,Le.asMonths=ze,Le.asYears=Ae,Le.valueOf=Kc,Le._bubble=Gc,Le.get=Mc,Le.milliseconds=Be,Le.seconds=Ce,Le.minutes=De,Le.hours=Ee,Le.days=Fe,Le.weeks=Oc,Le.months=Ge,Le.years=He,Le.humanize=Sc,Le.toISOString=Tc,Le.toString=Tc,Le.toJSON=Tc,Le.locale=ub,Le.localeData=vb,Le.toIsoString=fa("toIsoString() is deprecated. Please use toISOString() instead (notice the capitals)",Tc),Le.lang=Zd,J("X",0,0,"unix"),J("x",0,0,"valueOf"),O("x",od),O("X",rd),S("X",function(a,b,c){c._d=new Date(1e3*parseFloat(a,10))}),S("x",function(a,b,c){c._d=new Date(r(a))}),a.version="2.11.1",b(Ea),a.fn=ke,a.min=Ga,a.max=Ha,a.now=Td,a.utc=h,a.unix=kc,a.months=wc,a.isDate=d,a.locale=x,a.invalid=l,a.duration=Za,a.isMoment=p,a.weekdays=yc,a.parseZone=lc,a.localeData=z,a.isDuration=Ja,a.monthsShort=xc,a.weekdaysMin=Ac,a.defineLocale=y,a.weekdaysShort=zc,a.normalizeUnits=B,a.relativeTimeThreshold=Rc,a.prototype=ke;var Me=a;return Me});

// Datetime picker
// https://github.com/T00rk/bootstrap-material-datetimepicker
//
//

(function ($, moment) {
  var pluginName = 'materialDatePicker';
  var pluginDataName = 'plugin_' + pluginName;

  moment.locale('en');

  function Plugin(element, options) {
    this.currentView = 0;

    this.minDate = null;
    this.maxDate = null;

    this._attachedEvents = [];

    this.element = element;
    this.$element = $(element);

    this.params = {
      date: true,
      time: true,
      format: 'YYYY-MM-DD',
      minDate: null,
      maxDate: null,
      currentDate: null,
      lang: 'en',
      weekStart: 0,
      shortTime: false,
      cancelText: 'Cancel',
      okText: 'OK'
    };
    this.params = $.fn.extend(this.params, options);

    this.name = 'qor_' + this.setName();
    this.$element.attr('data-dtp', this.name);

    this.init();
  }

  $.fn[pluginName] = function (options, p) {
    this.each(function () {
      if (!$.data(this, pluginDataName)) {
        $.data(this, pluginDataName, new Plugin(this, options));
      } else {
        if (typeof ($.data(this, pluginDataName)[options]) === 'function') {
          $.data(this, pluginDataName)[options](p);
        }
        if (options === 'destroy') {
          delete $.data(this, pluginDataName);
        }
      }
    });
    return this;
  };

  Plugin.prototype = {
    init: function () {
      this.initDays();
      this.initDates();

      this.initTemplate();

      this.initButtons();

      this._attachEvent($(window), 'resize', this._centerBox(this));
      this._attachEvent(this.$dtpElement.find('.qor-datepicker__content'), 'click', this._onElementClick.bind(this));
      this._attachEvent(this.$dtpElement, 'click', this._onBackgroundClick.bind(this));
      this._attachEvent(this.$dtpElement.find('.qor-close > a'), 'click', this._onCloseClick.bind(this));
      this._attachEvent(this.$element, 'focus', this._onFocus.bind(this));
    },
    initDays: function () {
      this.days = [];
      for (var i = this.params.weekStart; this.days.length < 7; i++) {
        if (i > 6) {
          i = 0;
        }
        this.days.push(i.toString());
      }
    },
    initDates: function () {
      if (this.$element.val().length > 0) {
        if (typeof (this.params.format) !== 'undefined' && this.params.format !== null) {
          this.currentDate = moment(this.$element.val(), this.params.format).locale(this.params.lang);
        } else {
          this.currentDate = moment(this.$element.val()).locale(this.params.lang);
        }
      } else {
        if (typeof (this.$element.attr('value')) !== 'undefined' && this.$element.attr('value') !== null && this.$element.attr('value') !== '') {
          if (typeof (this.$element.attr('value')) === 'string') {
            if (typeof (this.params.format) !== 'undefined' && this.params.format !== null) {
              this.currentDate = moment(this.$element.attr('value'), this.params.format).locale(this.params.lang);
            } else {
              this.currentDate = moment(this.$element.attr('value')).locale(this.params.lang);
            }
          }
        } else {
          if (typeof (this.params.currentDate) !== 'undefined' && this.params.currentDate !== null) {
            if (typeof (this.params.currentDate) === 'string') {
              if (typeof (this.params.format) !== 'undefined' && this.params.format !== null) {
                this.currentDate = moment(this.params.currentDate, this.params.format).locale(this.params.lang);
              } else {
                this.currentDate = moment(this.params.currentDate).locale(this.params.lang);
              }
            } else {
              if (typeof (this.params.currentDate.isValid) === 'undefined' || typeof (this.params.currentDate.isValid) !== 'function') {
                var x = this.params.currentDate.getTime();
                this.currentDate = moment(x, 'x').locale(this.params.lang);
              } else {
                this.currentDate = this.params.currentDate;
              }
            }
            this.$element.val(this.currentDate.format(this.params.format));
          } else
            this.currentDate = moment();
        }
      }

      if (typeof (this.params.minDate) !== 'undefined' && this.params.minDate !== null) {
        if (typeof (this.params.minDate) === 'string') {
          if (typeof (this.params.format) !== 'undefined' && this.params.format !== null) {
            this.minDate = moment(this.params.minDate, this.params.format).locale(this.params.lang);
          } else {
            this.minDate = moment(this.params.minDate).locale(this.params.lang);
          }
        } else {
          if (typeof (this.params.minDate.isValid) === 'undefined' || typeof (this.params.minDate.isValid) !== 'function') {
            var x = this.params.minDate.getTime();
            this.minDate = moment(x, 'x').locale(this.params.lang);
          } else {
            this.minDate = this.params.minDate;
          }
        }
      }

      if (typeof (this.params.maxDate) !== 'undefined' && this.params.maxDate !== null) {
        if (typeof (this.params.maxDate) === 'string') {
          if (typeof (this.params.format) !== 'undefined' && this.params.format !== null) {
            this.maxDate = moment(this.params.maxDate, this.params.format).locale(this.params.lang);
          } else {
            this.maxDate = moment(this.params.maxDate).locale(this.params.lang);
          }
        } else {
          if (typeof (this.params.maxDate.isValid) === 'undefined' || typeof (this.params.maxDate.isValid) !== 'function') {
            var x = this.params.maxDate.getTime();
            this.maxDate = moment(x, 'x').locale(this.params.lang);
          } else {
            this.maxDate = this.params.maxDate;
          }
        }
      }

      if (!this.isAfterMinDate(this.currentDate)) {
        this.currentDate = moment(this.minDate);
      }
      if (!this.isBeforeMaxDate(this.currentDate)) {
        this.currentDate = moment(this.maxDate);
      }
    },
    initTemplate: function () {
      this.template = '<div class="qor-datepicker hidden" id="' + this.name + '">' +
        '<div class="qor-datepicker__content">' +
        '<div class="qor-date-view">' +
        '<header class="qor-header">' +
        '<div class="qor-actual-day">Lundi</div>' +
        '<div class="qor-close"><a href="javascript:void(0);"><i class="material-icons">clear</i></</div>' +
        '</header>' +

        '<div class="qor-date hidden">' +
        '<div>' +
        '<div class="left center p10">' +
        '<a href="javascript:void(0);" class="qor-select-year-before"><i class="material-icons">chevron_left</i></a>' +
        '</div>' +
        '<span class="qor-actual-year">2014</span>' +
        '<div class="right center p10">' +
        '<a href="javascript:void(0);" class="qor-select-year-after"><i class="material-icons">chevron_right</i></a>' +
        '</div>' +
        '</div>' +

        '<div>' +
        '<div class="left center p10">' +
        '<a href="javascript:void(0);" class="qor-select-month-before"><i class="material-icons">chevron_left</i></a>' +
        '</div>' +
        '<span><span class="qor-actual-month ">MAR</span> <span class="qor-actual-num">13</span></span>' +
        '<div class="right center p10">' +
        '<a href="javascript:void(0);" class="qor-select-month-after"><i class="material-icons">chevron_right</i></a>' +
        '</div>' +
        '</div>' +

        '</div>' +

        '<div class="qor-time hidden">' +
        '<div class="qor-actual-maxtime">23:55</div>' +
        '</div>' +
        '<div class="qor-picker">' +
        '<div class="qor-picker-calendar"></div>' +
        '<div class="qor-picker-datetime hidden">' +
        '<div class="qor-actual-meridien">' +
        '<div class="left p20">' +
        '<a class="qor-meridien-am" href="javascript:void(0);">AM</a>' +
        '</div>' +
        '<div class="qor-actual-time p60"></div>' +
        '<div class="right p20">' +
        '<a class="qor-meridien-pm" href="javascript:void(0);">PM</a>' +
        '</div>' +
        '</div>' +
        '<div class="qor-picker-clock"></div>' +
        '</div>' +
        '</div>' +
        '</div>' +
        '<div class="qor-datepicker__buttons">' +
        '<button class="qor-datepicker__save mdl-button mdl-button--colored mdl-js-button mdl-js-ripple-effect">' + this.params.okText + '</button>' +
        '<button class="qor-datepicker__cancel mdl-button mdl-button--colored mdl-js-button mdl-js-ripple-effect">' + this.params.cancelText + '</button>' +
        '</div>' +
        '</div>' +
        '</div>';

      if ($('body').find('#' + this.name).length <= 0) {
        $('body').append(this.template);

        this.dtpElement = $('body').find('#' + this.name);
        this.$dtpElement = $(this.dtpElement);
      }
    },
    initButtons: function () {
      this._attachEvent(this.$dtpElement.find('.qor-datepicker__cancel'), 'click', this._onCancelClick.bind(this));
      this._attachEvent(this.$dtpElement.find('.qor-datepicker__save'), 'click', this._onOKClick.bind(this));
      this._attachEvent(this.$dtpElement.find('a.qor-select-month-before'), 'click', this._onMonthBeforeClick.bind(this));
      this._attachEvent(this.$dtpElement.find('a.qor-select-month-after'), 'click', this._onMonthAfterClick.bind(this));
      this._attachEvent(this.$dtpElement.find('a.qor-select-year-before'), 'click', this._onYearBeforeClick.bind(this));
      this._attachEvent(this.$dtpElement.find('a.qor-select-year-after'), 'click', this._onYearAfterClick.bind(this));
    },
    initMeridienButtons: function () {
      this.$dtpElement.find('a.qor-meridien-am').off('click').on('click', this._onSelectAM.bind(this));
      this.$dtpElement.find('a.qor-meridien-pm').off('click').on('click', this._onSelectPM.bind(this));
    },
    initDate: function (d) {
      this.currentView = 0;

      this.$dtpElement.find('.qor-picker-calendar').removeClass('hidden');
      this.$dtpElement.find('.qor-picker-datetime').addClass('hidden');

      var _date = ((typeof (this.currentDate) !== 'undefined' && this.currentDate !== null) ? this.currentDate : null);
      var _calendar = this.generateCalendar(this.currentDate);

      if (typeof (_calendar.week) !== 'undefined' && typeof (_calendar.days) !== 'undefined') {
        var _template = this.constructHTMLCalendar(_date, _calendar);

        this.$dtpElement.find('a.qor-select-day').off('click');
        this.$dtpElement.find('.qor-picker-calendar').html(_template);

        this.$dtpElement.find('a.qor-select-day').on('click', this._onSelectDate.bind(this));

        this.toggleButtons(_date);
      }

      this._centerBox();
      this.showDate(_date);
    },
    initHours: function () {
      this.currentView = 1;

      if (!this.params.date) {
        var w = this.$dtpElement.find('.qor-datepicker__content').width();

        var ml = this.$dtpElement.find('.qor-picker-clock').css('marginLeft').replace('px', '');
        var mr = this.$dtpElement.find('.qor-picker-clock').css('marginRight').replace('px', '');

        var pl = this.$dtpElement.find('.qor-picker').css('paddingLeft').replace('px', '');
        var pr = this.$dtpElement.find('.qor-picker').css('paddingRight').replace('px', '');

        this.$dtpElement.find('.qor-picker-clock').innerWidth(w - (parseInt(ml) + parseInt(mr) + parseInt(pl) + parseInt(pr)));
      }

      this.showTime(this.currentDate);
      this.initMeridienButtons();

      this.$dtpElement.find('.qor-picker-datetime').removeClass('hidden');
      this.$dtpElement.find('.qor-picker-calendar').addClass('hidden');

      if (this.currentDate.hour() < 12) {
        this.$dtpElement.find('a.qor-meridien-am').click();
      } else {
        this.$dtpElement.find('a.qor-meridien-pm').click();
      }

      var hours = [];

      for (var h = 0; h < 12; ++h) {

        var hour = $('<div>', {
          class: 'qor-picker-time qor-picker-time-' + (h + 1)
        });

        var cH = ((this.currentDate.format('h') == 12) ? 0 : this.currentDate.format('h'));
        var hourLink = $('<a>', {
          href: 'javascript:void(0);',
          class: 'qor-select-hour'
        }).data('hour', h).text((h == 0 ? 12 : h));
        if (h == parseInt(cH))
          hourLink.addClass('selected');

        hour.append(hourLink);
        hours.push(hour);
      }

      this.$dtpElement.find('a.qor-select-hour').off('click');

      this.$dtpElement.find('.qor-picker-clock').html(hours);
      this.toggleTime(true);

      this.initHands(true);
    },
    initMinutes: function () {
      this.currentView = 2;

      this.showTime(this.currentDate);

      this.initMeridienButtons();

      if (this.currentDate.hour() < 12) {
        this.$dtpElement.find('a.qor-meridien-am').click();
      } else {
        this.$dtpElement.find('a.qor-meridien-pm').click();
      }

      this.$dtpElement.find('.qor-picker-calendar').addClass('hidden');
      this.$dtpElement.find('.qor-picker-datetime').removeClass('hidden');

      var minutes = [];

      for (var m = 0; m < 60; m += 5) {

        var minute = $('<div>', {
          class: 'qor-picker-time qor-picker-time-' + (m / 5 + 1)
        });

        var minuteLink = $('<a>', {
          href: 'javascript:void(0);',
          class: 'qor-select-minute'
        }).data('minute', m).text(((m.toString().length == 2) ? m : '0' + m));
        if (m == 5 * Math.round(this.currentDate.minute() / 5)) {
          minuteLink.addClass('selected');
          this.currentDate.minute(m);
        }

        minute.append(minuteLink);
        minutes.push(minute);
      }

      this.$dtpElement.find('a.qor-select-minute').off('click');

      this.$dtpElement.find('.qor-picker-clock').html(minutes);
      this.toggleTime(false);

      this.initHands(false);

      this._centerBox();
    },
    initHands: function (t) {
      this.$dtpElement.find('.qor-picker-clock').append(
        '<div class="qor-hand qor-hour-hand"></div>' +
        '<div class="qor-hand qor-minute-hand"></div>' +
        '<div class="qor-clock-center"></div>'
      );

      this.$dtpElement.find('.qor-hour-hand').addClass((t === true) ? 'on' : '');
      this.$dtpElement.find('.qor-minute-hand').addClass((t === false) ? 'on' : '');

      this.animateHands();

      this._centerBox();
    },
    animateHands: function () {
      var h = this.currentDate.hour();
      var m = this.currentDate.minute();

      this.rotateElement(this.$dtpElement.find('.qor-hour-hand'), (360 / 12) * h);
      this.rotateElement(this.$dtpElement.find('.qor-minute-hand'), ((360 / 60) * (5 * Math.round(this.currentDate.minute() / 5))));
    },
    isAfterMinDate: function (date, checkHour, checkMinute) {
      var _return = true;

      if (typeof (this.minDate) !== 'undefined' && this.minDate !== null) {
        var _minDate = moment(this.minDate);
        var _date = moment(date);

        if (!checkHour && !checkMinute) {
          _minDate.hour(0);
          _minDate.minute(0);

          _date.hour(0);
          _date.minute(0);
        }

        _minDate.second(0);
        _date.second(0);
        _minDate.millisecond(0);
        _date.millisecond(0);

        if (!checkMinute) {
          _date.minute(0);
          _minDate.minute(0);

          _return = (parseInt(_date.format('X')) >= parseInt(_minDate.format('X')));
        } else {
          _return = (parseInt(_date.format('X')) >= parseInt(_minDate.format('X')));
        }
      }

      return _return;
    },
    isBeforeMaxDate: function (date, checkTime, checkMinute) {
      var _return = true;

      if (typeof (this.maxDate) !== 'undefined' && this.maxDate !== null) {
        var _maxDate = moment(this.maxDate);
        var _date = moment(date);

        if (!checkTime && !checkMinute) {
          _maxDate.hour(0);
          _maxDate.minute(0);

          _date.hour(0);
          _date.minute(0);
        }

        _maxDate.second(0);
        _date.second(0);
        _maxDate.millisecond(0);
        _date.millisecond(0);

        if (!checkMinute) {
          _date.minute(0);
          _maxDate.minute(0);

          _return = (parseInt(_date.format('X')) <= parseInt(_maxDate.format('X')));
        } else {
          _return = (parseInt(_date.format('X')) <= parseInt(_maxDate.format('X')));
        }
      }

      return _return;
    },
    rotateElement: function (el, deg) {
      $(el).css({
        WebkitTransform: 'rotate(' + deg + 'deg)',
        '-moz-transform': 'rotate(' + deg + 'deg)'
      });
    },
    showDate: function (date) {
      if (date) {
        this.$dtpElement.find('.qor-actual-day').html(date.locale(this.params.lang).format('dddd'));
        this.$dtpElement.find('.qor-actual-month').html(date.locale(this.params.lang).format('MMM').toUpperCase());
        this.$dtpElement.find('.qor-actual-num').html(date.locale(this.params.lang).format('DD'));
        this.$dtpElement.find('.qor-actual-year').html(date.locale(this.params.lang).format('YYYY'));
      }
    },
    showTime: function (date) {
      if (date) {
        var minutes = (5 * Math.round(date.minute() / 5));
        var content = ((this.params.shortTime) ? date.format('hh') : date.format('HH')) + ':' + ((minutes.toString().length == 2) ? minutes : '0' + minutes);

        if (this.params.date)
          this.$dtpElement.find('.qor-actual-time').html(content);
        else {
          if (this.params.shortTime)
            this.$dtpElement.find('.qor-actual-day').html(date.format('A'));
          else
            this.$dtpElement.find('.qor-actual-day').html(' ');

          this.$dtpElement.find('.qor-actual-maxtime').html(content);
        }
      }
    },
    selectDate: function (date) {
      if (date) {
        this.currentDate.date(date);

        this.showDate(this.currentDate);
        this.$element.trigger('dateSelected', this.currentDate);
      }
    },
    generateCalendar: function (date) {
      var _calendar = {};

      if (date !== null) {
        var startOfMonth = moment(date).locale(this.params.lang).startOf('month');
        var endOfMonth = moment(date).locale(this.params.lang).endOf('month');

        var iNumDay = startOfMonth.format('d');

        _calendar.week = this.days;
        _calendar.days = [];

        for (var i = startOfMonth.date(); i <= endOfMonth.date(); i++) {
          if (i === startOfMonth.date()) {
            var iWeek = _calendar.week.indexOf(iNumDay.toString());
            if (iWeek > 0) {
              for (var x = 0; x < iWeek; x++) {
                _calendar.days.push(0);
              }
            }
          }
          _calendar.days.push(moment(startOfMonth).locale(this.params.lang).date(i));
        }
      }

      return _calendar;
    },
    constructHTMLCalendar: function (date, calendar) {
      var _template = '';

      _template += '<div class="qor-picker-month">' + date.locale(this.params.lang).format('MMMM YYYY') + '</div>';
      _template += '<table class="table qor-picker-days"><thead>';
      for (var i = 0; i < calendar.week.length; i++) {
        _template += '<th>' + moment(parseInt(calendar.week[i]), 'd').locale(this.params.lang).format('dd').substring(0, 1) + '</th>';
      }

      _template += '</thead>';
      _template += '<tbody><tr>';

      for (var i = 0; i < calendar.days.length; i++) {
        if (i % 7 == 0)
          _template += '</tr><tr>';
        _template += '<td data-date="' + moment(calendar.days[i]).locale(this.params.lang).format('D') + '">';
        if (calendar.days[i] != 0) {
          if (this.isBeforeMaxDate(moment(calendar.days[i]), false, false) === false || this.isAfterMinDate(moment(calendar.days[i]), false, false) === false) {
            _template += '<span class="qor-select-day">' + moment(calendar.days[i]).locale(this.params.lang).format('DD') + '</span>';
          } else {
            if (moment(calendar.days[i]).locale(this.params.lang).format('DD') === moment(this.currentDate).locale(this.params.lang).format('DD')) {
              _template += '<a href="javascript:void(0);" class="qor-select-day selected">' + moment(calendar.days[i]).locale(this.params.lang).format('DD') + '</a>';
            } else {
              _template += '<a href="javascript:void(0);" class="qor-select-day">' + moment(calendar.days[i]).locale(this.params.lang).format('DD') + '</a>';
            }
          }

          _template += '</td>';
        }
      }
      _template += '</tr></tbody></table>';

      return _template;
    },
    setName: function () {
      var text = '';
      var possible = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';

      for (var i = 0; i < 5; i++) {
        text += possible.charAt(Math.floor(Math.random() * possible.length));
      }

      return text;
    },
    isPM: function () {
      return this.$dtpElement.find('a.qor-meridien-pm').hasClass('selected');
    },
    setElementValue: function () {
      this.$element.trigger('beforeChange', this.currentDate);
      if (typeof ($.material) !== 'undefined') {
        this.$element.removeClass('empty');
      }
      this.$element.val(moment(this.currentDate).locale(this.params.lang).format(this.params.format));
      this.$element.trigger('change', this.currentDate);
    },
    toggleButtons: function (date) {
      if (date && date.isValid()) {
        var startOfMonth = moment(date).locale(this.params.lang).startOf('month');
        var endOfMonth = moment(date).locale(this.params.lang).endOf('month');

        if (!this.isAfterMinDate(startOfMonth, false, false)) {
          this.$dtpElement.find('a.qor-select-month-before').addClass('invisible');
        } else {
          this.$dtpElement.find('a.qor-select-month-before').removeClass('invisible');
        }

        if (!this.isBeforeMaxDate(endOfMonth, false, false)) {
          this.$dtpElement.find('a.qor-select-month-after').addClass('invisible');
        } else {
          this.$dtpElement.find('a.qor-select-month-after').removeClass('invisible');
        }

        var startOfYear = moment(date).locale(this.params.lang).startOf('year');
        var endOfYear = moment(date).locale(this.params.lang).endOf('year');

        if (!this.isAfterMinDate(startOfYear, false, false)) {
          this.$dtpElement.find('a.qor-select-year-before').addClass('invisible');
        } else {
          this.$dtpElement.find('a.qor-select-year-before').removeClass('invisible');
        }

        if (!this.isBeforeMaxDate(endOfYear, false, false)) {
          this.$dtpElement.find('a.qor-select-year-after').addClass('invisible');
        } else {
          this.$dtpElement.find('a.qor-select-year-after').removeClass('invisible');
        }
      }
    },
    toggleTime: function (isHours) {
      if (isHours) {
        this.$dtpElement.find('a.qor-select-hour').removeClass('disabled');
        this.$dtpElement.find('a.qor-select-hour').removeProp('disabled');
        this.$dtpElement.find('a.qor-select-hour').off('click');

        var _self = this;

        this.$dtpElement.find('a.qor-select-hour').each(function () {
          var _hour = $(this).data('hour');

          var _date = moment(_self.currentDate);
          _date.hour(_self.convertHours(_hour)).minute(0).second(0);

          if (_self.isAfterMinDate(_date, true, false) === false || _self.isBeforeMaxDate(_date, true, false) === false) {
            $(this).prop('disabled');
            $(this).addClass('disabled');
          } else {
            $(this).on('click', _self._onSelectHour.bind(_self));
          }
        });
      } else {
        this.$dtpElement.find('a.qor-select-minute').removeClass('disabled');
        this.$dtpElement.find('a.qor-select-minute').removeProp('disabled');
        this.$dtpElement.find('a.qor-select-minute').off('click');

        var _self = this;

        this.$dtpElement.find('a.qor-select-minute').each(function () {
          var _minute = $(this).data('minute');

          var _date = moment(_self.currentDate);
          _date.minute(_minute).second(0);

          if (_self.isAfterMinDate(_date, true, true) === false || _self.isBeforeMaxDate(_date, true, true) === false) {
            $(this).prop('disabled');
            $(this).addClass('disabled');
          } else {
            $(this).on('click', _self._onSelectMinute.bind(_self));
          }
        });
      }
    },
    _attachEvent: function (el, ev, fn) {
      el.on(ev, fn);
      this._attachedEvents.push([el, ev, fn]);
    },
    _detachEvents: function () {
      for (var i = this._attachedEvents.length - 1; i >= 0; i--) {
        this._attachedEvents[i][0].off(this._attachedEvents[i][1], this._attachedEvents[i][2]);
        this._attachedEvents.splice(i, 1);
      }
    },
    _onFocus: function () {
      this.currentView = 0;
      this.$element.blur();

      this.initDates();

      this.show();

      if (this.params.date) {
        this.$dtpElement.find('.qor-date').removeClass('hidden');
        this.initDate();
      } else {
        if (this.params.time) {
          this.$dtpElement.find('.qor-time').removeClass('hidden');
          this.initHours();
        }
      }
    },
    _onBackgroundClick: function (e) {
      e.stopPropagation();
      this.hide();
    },
    _onElementClick: function (e) {
      e.stopPropagation();
    },
    _onCloseClick: function () {
      this.hide();
    },
    _onOKClick: function () {
      switch (this.currentView) {
      case 0:
        if (this.params.time === true) {
          this.initHours();
        } else {
          this.setElementValue();
          this.hide();
        }
        break;
      case 1:
        this.initMinutes();
        break;
      case 2:
        this.setElementValue();
        this.hide();
        break;
      }
    },
    _onCancelClick: function () {
      if (this.params.time) {
        switch (this.currentView) {
        case 0:
          this.hide();
          break;
        case 1:
          if (this.params.date) {
            this.initDate();
          } else {
            this.hide();
          }
          break;
        case 2:
          this.initHours();
          break;
        }
      } else {
        this.hide();
      }
    },
    _onMonthBeforeClick: function () {
      this.currentDate.subtract(1, 'months');
      this.initDate(this.currentDate);
    },
    _onMonthAfterClick: function () {
      this.currentDate.add(1, 'months');
      this.initDate(this.currentDate);
    },
    _onYearBeforeClick: function () {
      this.currentDate.subtract(1, 'years');
      this.initDate(this.currentDate);
    },
    _onYearAfterClick: function () {
      this.currentDate.add(1, 'years');
      this.initDate(this.currentDate);
    },
    _onSelectDate: function (e) {
      this.$dtpElement.find('a.qor-select-day').removeClass('selected');
      $(e.currentTarget).addClass('selected');

      this.selectDate($(e.currentTarget).parent().data('date'));
    },
    _onSelectHour: function (e) {
      this.$dtpElement.find('a.qor-select-hour').removeClass('selected');
      $(e.currentTarget).addClass('selected');

      var dataHour = parseInt($(e.currentTarget).data('hour'));
      if (this.isPM())
        dataHour += 12;

      this.currentDate.hour(dataHour);
      this.showTime(this.currentDate);

      this.animateHands();
    },
    _onSelectMinute: function (e) {
      this.$dtpElement.find('a.qor-select-minute').removeClass('selected');
      $(e.currentTarget).addClass('selected');

      this.currentDate.minute(parseInt($(e.currentTarget).data('minute')));
      this.showTime(this.currentDate);

      this.animateHands();
    },
    _onSelectAM: function (e) {
      $('.qor-actual-meridien').find('a').removeClass('selected');
      $(e.currentTarget).addClass('selected');

      if (this.currentDate.hour() >= 12) {
        if (this.currentDate.subtract(12, 'hours'))
          this.showTime(this.currentDate);
      }
      this.toggleTime((this.currentView === 1));
    },
    _onSelectPM: function (e) {
      $('.qor-actual-meridien').find('a').removeClass('selected');
      $(e.currentTarget).addClass('selected');

      if (this.currentDate.hour() < 12) {
        if (this.currentDate.add(12, 'hours'))
          this.showTime(this.currentDate);
      }
      this.toggleTime((this.currentView === 1));
    },
    convertHours: function (h) {
      var _return = h;

      if ((h < 12) && this.isPM())
        _return += 12;

      return _return;
    },
    setDate: function (date) {
      this.params.currentDate = date;
      this.initDates();
    },
    setMinDate: function (date) {
      this.params.minDate = date;
      this.initDates();
    },
    setMaxDate: function (date) {
      this.params.maxDate = date;
      this.initDates();
    },
    destroy: function () {
      this._detachEvents();
      this.$dtpElement.remove();
    },
    show: function () {
      this.$dtpElement.removeClass('hidden');
      this._centerBox();
    },
    hide: function () {
      this.$dtpElement.addClass('hidden');
    },
    resetDate: function () {

    },
    _centerBox: function () {
      var h = (this.$dtpElement.height() - this.$dtpElement.find('.qor-datepicker__content').height()) / 2;
      this.$dtpElement.find('.qor-datepicker__content').css('marginLeft', -(this.$dtpElement.find('.qor-datepicker__content').width() / 2) + 'px');
      this.$dtpElement.find('.qor-datepicker__content').css('top', h + 'px');
    }
  };
})(jQuery, moment);
