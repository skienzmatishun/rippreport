var 
_____WB$wombat$assign$function_____=function(name){return(self._wb_wombat&&self._wb_wombat.local_init&&self._wb_wombat.local_init(name))||self[name];};if(!self.__WB_pmw){self.__WB_pmw=function(obj){this.__WB_source=obj;return 
this;}}
{let window=_____WB$wombat$assign$function_____("window");let 
self=_____WB$wombat$assign$function_____("self");let 
document=_____WB$wombat$assign$function_____("document");let 
location=_____WB$wombat$assign$function_____("location");let 
top=_____WB$wombat$assign$function_____("top");let 
parent=_____WB$wombat$assign$function_____("parent");let 
frames=_____WB$wombat$assign$function_____("frames");let 
opener=_____WB$wombat$assign$function_____("opener");parcelRequire=function(e,r,t,n){var 
i,o="function"==typeof parcelRequire&&parcelRequire,u="function"==typeof 
require&&require;function f(t,n){if(!r[t]){if(!e[t]){var 
i="function"==typeof parcelRequire&&parcelRequire;if(!n&&i)return 
i(t,!0);if(o)return o(t,!0);if(u&&"string"==typeof t)return u(t);var c=new 
Error("Cannot find module '"+t+"'");throw 
c.code="MODULE_NOT_FOUND",c}p.resolve=function(r){return 
e[t][1][r]||r},p.cache={};var l=r[t]=new 
f.Module(t);e[t][0].call(l.exports,p,l,l.exports,this)}return 
r[t].exports;function p(e){return 
f(p.resolve(e))}}f.isParcelRequire=!0,f.Module=function(e){this.id=e,this.bundle=f,this.exports={}},f.modules=e,f.cache=r,f.parent=o,f.register=function(r,t){e[r]=[function(e,r){r.exports=t},{}]};for(var 
c=0;c<t.length;c++)try{f(t[c])}catch(e){i||(i=e)}if(t.length){var 
l=f(t[t.length-1]);"object"==typeof exports&&"undefined"!=typeof 
module?module.exports=l:"function"==typeof 
define&&define.amd?define(function(){return 
l}):n&&(this[n]=l)}if(parcelRequire=f,i)throw i;return 
f}({"asWa":[function(require,module,exports){!function(r){"use 
strict";function n(r,n,t){return t.a=r,t.f=n,t}function t(r){return 
n(2,r,function(n){return function(t){return r(n,t)}})}function e(r){return 
n(3,r,function(n){return function(t){return function(e){return 
r(n,t,e)}}})}function u(r){return n(4,r,function(n){return 
function(t){return function(e){return function(u){return 
r(n,t,e,u)}}}})}function a(r){return n(5,r,function(n){return 
function(t){return function(e){return function(u){return 
function(a){return r(n,t,e,u,a)}}}}})}function c(r){return 
n(6,r,function(n){return function(t){return function(e){return 
function(u){return function(a){return function(c){return 
r(n,t,e,u,a,c)}}}}}})}function i(r){return n(7,r,function(n){return 
function(t){return function(e){return function(u){return 
function(a){return function(c){return function(i){return 
r(n,t,e,u,a,c,i)}}}}}}})}function o(r){return n(9,r,function(n){return 
function(t){return function(e){return function(u){return 
function(a){return function(c){return function(i){return 
function(o){return function(f){return 
r(n,t,e,u,a,c,i,o,f)}}}}}}}}})}function f(r,n,t){return 
2===r.a?r.f(n,t):r(n)(t)}function s(r,n,t,e){return 
3===r.a?r.f(n,t,e):r(n)(t)(e)}function l(r,n,t,e,u){return 
4===r.a?r.f(n,t,e,u):r(n)(t)(e)(u)}function b(r,n,t,e,u,a){return 
5===r.a?r.f(n,t,e,u,a):r(n)(t)(e)(u)(a)}function d(r,n,t,e,u,a,c){return 
6===r.a?r.f(n,t,e,u,a,c):r(n)(t)(e)(u)(a)(c)}function 
v(r,n,t,e,u,a,c,i){return 
7===r.a?r.f(n,t,e,u,a,c,i):r(n)(t)(e)(u)(a)(c)(i)}function p(r,n){for(var 
t,e=[],u=m(r,n,0,e);u&&(t=e.pop());u=m(t.a,t.b,0,e));return u}function 
m(r,n,t,e){if(r===n)return!0;if("object"!=typeof 
r||null===r||null===n)return"function"==typeof r&&B(5),!1;if(t>100)return 
e.push(y(r,n)),!0;for(var u in 
0>r.$&&(r=st(r),n=st(n)),r)if(!m(r[u],n[u],t+1,e))return!1;return!0}var 
g=t(p);function h(r,n,t){if("object"!=typeof r)return 
r===n?0:n>r?-1:1;if(void 
0===r.$)return(t=h(r.a,n.a))?t:(t=h(r.b,n.b))?t:h(r.c,n.c);for(;r.b&&n.b&&!(t=h(r.a,n.a));r=r.b,n=n.b);return 
t||(r.b?1:n.b?-1:0)}var $=t(function(r,n){var t=h(r,n);return 
0>t?it:t?ct:at}),w=0;function y(r,n){return{a:r,b:n}}function 
q(r,n,t){return{a:r,b:n,c:t}}function x(r){return r}function A(r,n){var 
t={};for(var e in r)t[e]=r[e];for(var e in n)t[e]=n[e];return t}var 
k=t(E);function E(r,n){if("string"==typeof r)return r+n;if(!r.b)return 
n;var t=D(r.a,n);r=r.b;for(var e=t;r.b;r=r.b)e=e.b=D(r.a,n);return t}var 
L={$:0};function D(r,n){return{$:1,a:r,b:n}}var T=t(D);function 
N(r){for(var n=L,t=r.length;t--;)n=D(r[t],n);return n}function 
S(r){for(var n=[];r.b;r=r.b)n.push(r.a);return n}var 
C=e(function(r,n,t){for(var 
e=[];n.b&&t.b;n=n.b,t=t.b)e.push(f(r,n.a,t.a));return 
N(e)}),V=t(function(r,n){return N(S(n).sort(function(n,t){return 
h(r(n),r(t))}))}),R=e(function(r,n,t){for(var 
e=[],u=0;r>u;u++)e[u]=t(n+u);return e}),U=t(function(r,n){for(var 
t=[],e=0;r>e&&n.b;e++)t[e]=n.a,n=n.b;return t.length=e,y(t,n)});function 
B(r){throw 
Error("https://github.com/elm/core/blob/1.0.0/hints/"+r+".md")}var 
O=t(function(r,n){return r+n}),P=t(function(r,n){return 
r*n}),j=t(Math.pow),I=t(function(r,n){var t=n%r;return 
0===r?B(11):t>0&&0>r||0>t&&r>0?t+r:t}),_=Math.ceil,G=Math.floor,z=Math.round,H=Math.log,M=t(function(r,n){return 
r+n}),F=e(function(r,n,t){for(var e=t.length,u=0;e>u;){var 
a=t[u],c=t.charCodeAt(u);u++,55296>c||c>56319||(a+=t[u],u++),n=f(r,x(a),n)}return 
n}),K=e(function(r,n,t){for(var e=t.length;e--;){var 
u=t[e],a=t.charCodeAt(e);56320>a||a>57343||(u=t[--e]+u),n=f(r,x(u),n)}return 
n}),Z=t(function(r,n){return n.split(r)}),J=t(function(r,n){return 
n.join(r)}),X=e(function(r,n,t){return 
t.slice(r,n)}),Y=t(function(r,n){for(var t=n.length;t--;){var 
e=n[t],u=n.charCodeAt(t);if(56320>u||u>57343||(e=n[--t]+e),!r(x(e)))return!1}return!0}),Q=t(function(r,n){return 
n.indexOf(r)>-1}),W=t(function(r,n){return 
0===n.indexOf(r)}),rr=t(function(r,n){return 
n.length>=r.length&&n.lastIndexOf(r)===n.length-r.length});function 
nr(r){return{$:2,b:r}}var tr=nr(function(r){return"number"!=typeof 
r?yr("an 
INT",r):r>-2147483647&&2147483647>r&&(0|r)===r?pt(r):!isFinite(r)||r%1?yr("an 
INT",r):pt(r)}),er=nr(function(r){return"boolean"==typeof r?pt(r):yr("a 
BOOL",r)}),ur=nr(function(r){return"number"==typeof r?pt(r):yr("a 
FLOAT",r)}),ar=nr(function(r){return 
pt(kr(r))}),cr=nr(function(r){return"string"==typeof r?pt(r):r instanceof 
String?pt(r+""):yr("a 
STRING",r)}),ir=t(function(r,n){return{$:6,d:r,b:n}});function 
or(r,n){return{$:9,f:r,g:n}}var 
fr=t(function(r,n){return{$:10,b:n,h:r}}),sr=t(function(r,n){return 
or(r,[n])}),lr=e(function(r,n,t){return 
or(r,[n,t])}),br=u(function(r,n,t,e){return 
or(r,[n,t,e])}),dr=a(function(r,n,t,e,u){return 
or(r,[n,t,e,u])}),vr=c(function(r,n,t,e,u,a){return 
or(r,[n,t,e,u,a])}),pr=t(function(r,n){try{return 
gr(r,JSON.parse(n))}catch(r){return lt(f(bt,"This is not valid JSON! 
"+r.message,kr(n)))}}),mr=t(function(r,n){return gr(r,Er(n))});function 
gr(r,n){switch(r.$){case 2:return r.b(n);case 5:return 
null===n?pt(r.c):yr("null",n);case 3:return $r(n)?hr(r.b,n,N):yr("a 
LIST",n);case 4:return $r(n)?hr(r.b,n,wr):yr("an ARRAY",n);case 6:var 
t=r.d;if("object"!=typeof n||null===n||!(t in n))return yr("an OBJECT with 
a field named `"+t+"`",n);var e=gr(r.b,n[t]);return 
se(e)?e:lt(f(dt,t,e.a));case 7:var u=r.e;return 
$r(n)?n.length>u?(e=gr(r.b,n[u]),se(e)?e:lt(f(vt,u,e.a))):yr("a LONGER 
array. Need index "+u+" but only see "+n.length+" entries",n):yr("an 
ARRAY",n);case 8:if("object"!=typeof n||null===n||$r(n))return yr("an 
OBJECT",n);var a=L;for(var c in 
n)if(n.hasOwnProperty(c)){if(e=gr(r.b,n[c]),!se(e))return 
lt(f(dt,c,e.a));a=D(y(c,e.a),a)}return pt(jt(a));case 9:for(var 
i=r.f,o=r.g,s=0;o.length>s;s++){if(e=gr(o[s],n),!se(e))return 
e;i=i(e.a)}return pt(i);case 10:return 
e=gr(r.b,n),se(e)?gr(r.h(e.a),n):e;case 11:for(var 
l=L,b=r.g;b.b;b=b.b){if(e=gr(b.a,n),se(e))return e;l=D(e.a,l)}return 
lt(mt(jt(l)));case 1:return lt(f(bt,r.a,kr(n)));case 0:return 
pt(r.a)}}function hr(r,n,t){for(var e=n.length,u=[],a=0;e>a;a++){var 
c=gr(r,n[a]);if(!se(c))return lt(f(vt,a,c.a));u[a]=c.a}return 
pt(t(u))}function $r(r){return Array.isArray(r)||"undefined"!=typeof 
FileList&&r instanceof FileList}function wr(r){return 
f(fe,r.length,function(n){return r[n]})}function yr(r,n){return 
lt(f(bt,"Expecting "+r,kr(n)))}function 
qr(r,n){if(r===n)return!0;if(r.$!==n.$)return!1;switch(r.$){case 0:case 
1:return r.a===n.a;case 2:return r.b===n.b;case 5:return r.c===n.c;case 
3:case 4:case 8:return qr(r.b,n.b);case 6:return 
r.d===n.d&&qr(r.b,n.b);case 7:return r.e===n.e&&qr(r.b,n.b);case 9:return 
r.f===n.f&&xr(r.g,n.g);case 10:return r.h===n.h&&qr(r.b,n.b);case 
11:return xr(r.g,n.g)}}function xr(r,n){var 
t=r.length;if(t!==n.length)return!1;for(var 
e=0;t>e;e++)if(!qr(r[e],n[e]))return!1;return!0}var 
Ar=t(function(r,n){return JSON.stringify(Er(n),null,r)+""});function 
kr(r){return r}function Er(r){return r}var Lr=e(function(r,n,t){return 
t[r]=Er(n),t});function Dr(r){return{$:0,a:r}}function 
Tr(r){return{$:2,b:r,c:null}}kr(null);var 
Nr=t(function(r,n){return{$:3,b:r,d:n}}),Sr=t(function(r,n){return{$:4,b:r,d:n}}),Cr=0;function 
Vr(r){var n={$:0,e:Cr++,f:r,g:null,h:[]};return jr(n),n}function 
Rr(r){return Tr(function(n){n(Dr(Vr(r)))})}function 
Ur(r,n){r.h.push(n),jr(r)}var Br=t(function(r,n){return 
Tr(function(t){Ur(r,n),t(Dr(w))})}),Or=!1,Pr=[];function 
jr(r){if(Pr.push(r),!Or){for(Or=!0;r=Pr.shift();)Ir(r);Or=!1}}function 
Ir(r){for(;r.f;){var 
n=r.f.$;if(0===n||1===n){for(;r.g&&r.g.$!==n;)r.g=r.g.i;if(!r.g)return;r.f=r.g.b(r.f.a),r.g=r.g.i}else{if(2===n)return 
void(r.f.c=r.f.b(function(n){r.f=n,jr(r)}));if(5===n){if(0===r.h.length)return;r.f=r.f.b(r.h.shift())}else 
r.g={$:3===n?0:1,b:r.f.b,i:r.g},r.f=r.f.d}}}function _r(r){return 
Tr(function(n){var t=setTimeout(function(){n(Dr(w))},r);return 
function(){clearTimeout(t)}})}var Gr={};function 
zr(r,n,t,e,u){return{b:r,c:n,d:t,e:e,f:u}}function Hr(r,n){var 
t={g:n,h:void 0},e=r.c,u=r.d,a=r.e,c=r.f;function i(r){return 
f(Nr,i,{$:5,b:function(n){var i=n.a;return 
0===n.$?s(u,t,i,r):a&&c?l(e,t,i.i,i.j,r):s(e,t,a?i.i:i.j,r)}})}return 
t.h=Vr(f(Nr,i,r.b))}var Mr=t(function(r,n){return 
Tr(function(t){r.g(n),t(Dr(w))})}),Fr=t(function(r,n){return 
f(Br,r.h,{$:0,a:n})});function Kr(r){return 
function(n){return{$:1,k:r,l:n}}}function Zr(r){return{$:2,m:r}}var 
Jr=t(function(r,n){return{$:3,n:r,o:n}}),Xr=[],Yr=!1;function 
Qr(r,n,t){if(Xr.push({p:r,q:n,r:t}),!Yr){Yr=!0;for(var 
e;e=Xr.shift();)Wr(e.p,e.q,e.r);Yr=!1}}function Wr(r,n,t){var e={};for(var 
u in 
rn(!0,n,e,null),rn(!1,t,e,null),r)Ur(r[u],{$:"fx",a:e[u]||{i:L,j:L}})}function 
rn(r,n,t,e){switch(n.$){case 1:var u=n.k,a=function(r,t,e){function 
u(r){for(var n=e;n;n=n.t)r=n.s(r);return r}return 
f(r?Gr[t].e:Gr[t].f,u,n.l)}(r,u,e);return void(t[u]=function(r,n,t){return 
t=t||{i:L,j:L},r?t.i=D(n,t.i):t.j=D(n,t.j),t}(r,a,t[u]));case 2:for(var 
c=n.m;c.b;c=c.b)rn(r,c.a,t,e);return;case 3:return void 
rn(r,n.o,t,{s:n.n,t:e})}}var nn,tn=t(function(r,n){return n});var 
en="undefined"!=typeof document?document:{};function 
un(r,n){r.appendChild(n)}function an(r){return{$:0,a:r}}var 
cn=t(function(r,n){return t(function(t,e){for(var u=[],a=0;e.b;e=e.b){var 
c=e.a;a+=c.b||0,u.push(c)}return 
a+=u.length,{$:1,c:n,d:$n(t),e:u,f:r,b:a}})}),on=cn(void 
0);t(function(r,n){return t(function(t,e){for(var u=[],a=0;e.b;e=e.b){var 
c=e.a;a+=c.b.b||0,u.push(c)}return 
a+=u.length,{$:2,c:n,d:$n(t),e:u,f:r,b:a}})})(void 0);var 
fn=t(function(r,n){return{$:4,j:r,k:n,b:1+(n.b||0)}}),sn=t(function(r,n){return{$:"a0",n:r,o:n}}),ln=t(function(r,n){return{$:"a1",n:r,o:n}}),bn=t(function(r,n){return{$:"a2",n:r,o:n}}),dn=t(function(r,n){return{$:"a3",n:r,o:n}});function 
vn(r){return/^\s*(javascript:|data:text\/html)/i.test(r)?"":r}var 
pn,mn=t(function(r,n){return"a0"===n.$?f(sn,n.n,function(r,n){var 
t=ve(n);return{$:n.$,a:t?s(be,3>t?gn:hn,de(r),n.a):f(le,r,n.a)}}(r,n.o)):n}),gn=t(function(r,n){return 
y(r(n.a),n.b)}),hn=t(function(r,n){return{aF:r(n.aF),bn:n.bn,bb:n.bb}});function 
$n(r){for(var n={};r.b;r=r.b){var t=r.a,e=t.$,u=t.n,a=t.o;if("a2"!==e){var 
c=n[e]||(n[e]={});"a3"===e&&"class"===u?wn(c,u,a):c[u]=a}else"className"===u?wn(n,u,Er(a)):n[u]=Er(a)}return 
n}function wn(r,n,t){var e=r[n];r[n]=e?e+" "+t:t}function yn(r,n){var 
t=r.$;if(5===t)return yn(r.k||(r.k=r.m()),n);if(0===t)return 
en.createTextNode(r.a);if(4===t){for(var 
e=r.k,u=r.j;4===e.$;)"object"!=typeof u?u=[u,e.j]:u.push(e.j),e=e.k;var 
a={j:u,p:n};return(c=yn(e,a)).elm_event_node_ref=a,c}if(3===t)return 
qn(c=r.h(r.g),n,r.d),c;var 
c=r.f?en.createElementNS(r.f,r.c):en.createElement(r.c);nn&&"a"==r.c&&c.addEventListener("click",nn(c)),qn(c,n,r.d);for(var 
i=r.e,o=0;i.length>o;o++)un(c,yn(1===t?i[o]:i[o].b,n));return c}function 
qn(r,n,t){for(var e in t){var 
u=t[e];"a1"===e?xn(r,u):"a0"===e?En(r,n,u):"a3"===e?An(r,u):"a4"===e?kn(r,u):("value"!==e&&"checked"!==e||r[e]!==u)&&(r[e]=u)}}function 
xn(r,n){var t=r.style;for(var e in n)t[e]=n[e]}function An(r,n){for(var t 
in n){var e=n[t];void 
0!==e?r.setAttribute(t,e):r.removeAttribute(t)}}function kn(r,n){for(var t 
in n){var e=n[t],u=e.f,a=e.o;void 
0!==a?r.setAttributeNS(u,t,a):r.removeAttributeNS(u,t)}}function 
En(r,n,t){var e=r.elmFs||(r.elmFs={});for(var u in t){var 
a=t[u],c=e[u];if(a){if(c){if(c.q.$===a.$){c.q=a;continue}r.removeEventListener(u,c)}c=Ln(n,a),r.addEventListener(u,c,pn&&{passive:2>ve(a)}),e[u]=c}else 
r.removeEventListener(u,c),e[u]=void 
0}}try{window.addEventListener("t",null,Object.defineProperty({},"passive",{get:function(){pn=!0}}))}catch(r){}function 
Ln(r,n){function t(n){var e=t.q,u=gr(e.a,n);if(se(u)){for(var 
a,c=ve(e),i=u.a,o=c?3>c?i.a:i.aF:i,f=1==c?i.b:3==c&&i.bn,s=(f&&n.stopPropagation(),(2==c?i.b:3==c&&i.bb)&&n.preventDefault(),r);a=s.j;){if("function"==typeof 
a)o=a(o);else for(var l=a.length;l--;)o=a[l](o);s=s.p}s(o,f)}}return 
t.q=n,t}function Dn(r,n){return r.$==n.$&&qr(r.a,n.a)}function 
Tn(r,n,t,e){var u={$:n,r:t,s:e,t:void 0,u:void 0};return 
r.push(u),u}function Nn(r,n,t,e){if(r!==n){var 
u=r.$,a=n.$;if(u!==a){if(1!==u||2!==a)return void 
Tn(t,0,e,n);n=function(r){for(var 
n=r.e,t=n.length,e=[],u=0;t>u;u++)e[u]=n[u].b;return{$:1,c:r.c,d:r.d,e:e,f:r.f,b:r.b}}(n),a=1}switch(a){case 
5:for(var 
c=r.l,i=n.l,o=c.length,f=o===i.length;f&&o--;)f=c[o]===i[o];if(f)return 
void(n.k=r.k);n.k=n.m();var s=[];return 
Nn(r.k,n.k,s,0),void(s.length>0&&Tn(t,1,e,s));case 4:for(var 
l=r.j,b=n.j,d=!1,v=r.k;4===v.$;)d=!0,"object"!=typeof 
l?l=[l,v.j]:l.push(v.j),v=v.k;for(var p=n.k;4===p.$;)d=!0,"object"!=typeof 
b?b=[b,p.j]:b.push(p.j),p=p.k;return d&&l.length!==b.length?void 
Tn(t,0,e,n):((d?function(r,n){for(var 
t=0;r.length>t;t++)if(r[t]!==n[t])return!1;return!0}(l,b):l===b)||Tn(t,2,e,b),void 
Nn(v,p,t,e+1));case 0:return void(r.a!==n.a&&Tn(t,3,e,n.a));case 1:return 
void Sn(r,n,t,e,Vn);case 2:return void Sn(r,n,t,e,Rn);case 
3:if(r.h!==n.h)return void Tn(t,0,e,n);var 
m=Cn(r.d,n.d);m&&Tn(t,4,e,m);var g=n.i(r.g,n.g);return 
void(g&&Tn(t,5,e,g))}}}function Sn(r,n,t,e,u){if(r.c===n.c&&r.f===n.f){var 
a=Cn(r.d,n.d);a&&Tn(t,4,e,a),u(r,n,t,e)}else Tn(t,0,e,n)}function 
Cn(r,n,t){var e;for(var u in 
r)if("a1"!==u&&"a0"!==u&&"a3"!==u&&"a4"!==u)if(u in n){var 
a=r[u],c=n[u];a===c&&"value"!==u&&"checked"!==u||"a0"===t&&Dn(a,c)||((e=e||{})[u]=c)}else(e=e||{})[u]=t?"a1"===t?"":"a0"===t||"a3"===t?void 
0:{f:r[u].f,o:void 0}:"string"==typeof r[u]?"":null;else{var 
i=Cn(r[u],n[u]||{},u);i&&((e=e||{})[u]=i)}for(var o in n)o in 
r||((e=e||{})[o]=n[o]);return e}function Vn(r,n,t,e){var 
u=r.e,a=n.e,c=u.length,i=a.length;c>i?Tn(t,6,e,{v:i,i:c-i}):i>c&&Tn(t,7,e,{v:c,e:a});for(var 
o=i>c?c:i,f=0;o>f;f++){var s=u[f];Nn(s,a[f],t,++e),e+=s.b||0}}function 
Rn(r,n,t,e){for(var 
u=[],a={},c=[],i=r.e,o=n.e,f=i.length,s=o.length,l=0,b=0,d=e;f>l&&s>b;){var 
v=(E=i[l]).a,p=(L=o[b]).a,m=E.b,g=L.b,h=void 0,$=void 0;if(v!==p){var 
w=i[l+1],y=o[b+1];if(w){var q=w.a,x=w.b;$=p===q}if(y){var 
A=y.a,k=y.b;h=v===A}if(h&&$)Nn(m,k,u,++d),Bn(a,u,v,g,b,c),d+=m.b||0,On(a,u,v,x,++d),d+=x.b||0,l+=2,b+=2;else 
if(h)d++,Bn(a,u,p,g,b,c),Nn(m,k,u,d),d+=m.b||0,l+=1,b+=2;else 
if($)On(a,u,v,m,++d),d+=m.b||0,Nn(x,g,u,++d),d+=x.b||0,l+=2,b+=1;else{if(!w||q!==A)break;On(a,u,v,m,++d),Bn(a,u,p,g,b,c),d+=m.b||0,Nn(x,k,u,++d),d+=x.b||0,l+=2,b+=2}}else 
Nn(m,g,u,++d),d+=m.b||0,l++,b++}for(;f>l;){var 
E;On(a,u,(E=i[l]).a,m=E.b,++d),d+=m.b||0,l++}for(;s>b;){var 
L,D=D||[];Bn(a,u,(L=o[b]).a,L.b,void 
0,D),b++}(u.length>0||c.length>0||D)&&Tn(t,8,e,{w:u,x:c,y:D})}var 
Un="_elmW6BL";function Bn(r,n,t,e,u,a){var c=r[t];if(!c)return 
a.push({r:u,A:c={c:0,z:e,r:u,s:void 
0}}),void(r[t]=c);if(1===c.c){a.push({r:u,A:c}),c.c=2;var i=[];return 
Nn(c.z,e,i,c.r),c.r=u,void(c.s.s={w:i,A:c})}Bn(r,n,t+Un,e,u,a)}function 
On(r,n,t,e,u){var a=r[t];if(a){if(0===a.c){a.c=2;var c=[];return 
Nn(e,a.z,c,u),void Tn(n,9,u,{w:c,A:a})}On(r,n,t+Un,e,u)}else{var 
i=Tn(n,9,u,void 0);r[t]={c:1,z:e,r:u,s:i}}}function Pn(r,n,t,e){return 
0===t.length?r:(function r(n,t,e,u){!function n(t,e,u,a,c,i,o){for(var 
f=u[a],s=f.r;s===c;){var l=f.$;if(1===l)r(t,e.k,f.s,o);else 
if(8===l)f.t=t,f.u=o,(b=f.s.w).length>0&&n(t,e,b,0,c,i,o);else 
if(9===l){f.t=t,f.u=o;var 
b,d=f.s;d&&(d.A.s=t,(b=d.w).length>0&&n(t,e,b,0,c,i,o))}else 
f.t=t,f.u=o;if(!(f=u[++a])||(s=f.r)>i)return a}var v=e.$;if(4===v){for(var 
p=e.k;4===p.$;)p=p.k;return n(t,p,u,a,c+1,i,t.elm_event_node_ref)}for(var 
m=e.e,g=t.childNodes,h=0;m.length>h;h++){var 
$=1===v?m[h]:m[h].b,w=++c+($.b||0);if(!(c>s||s>w||(f=u[a=n(g[h],$,u,a,c,w,o)])&&(s=f.r)<=i))return 
a;c=w}return a}(n,t,e,0,0,t.b,u)}(r,n,t,e),jn(r,t))}function 
jn(r,n){for(var t=0;n.length>t;t++){var 
e=n[t],u=e.t,a=In(u,e);u===r&&(r=a)}return r}function 
In(r,n){switch(n.$){case 0:return function(r){var 
t=r.parentNode,e=yn(n.s,n.u);return 
e.elm_event_node_ref||(e.elm_event_node_ref=r.elm_event_node_ref),t&&e!==r&&t.replaceChild(e,r),e}(r);case 
4:return qn(r,n.u,n.s),r;case 3:return 
r.replaceData(0,r.length,n.s),r;case 1:return jn(r,n.s);case 2:return 
r.elm_event_node_ref?r.elm_event_node_ref.j=n.s:r.elm_event_node_ref={j:n.s,p:n.u},r;case 
6:for(var t=n.s,e=0;t.i>e;e++)r.removeChild(r.childNodes[t.v]);return 
r;case 7:for(var 
u=(t=n.s).e,a=r.childNodes[e=t.v];u.length>e;e++)r.insertBefore(yn(u[e],n.u),a);return 
r;case 9:if(!(t=n.s))return r.parentNode.removeChild(r),r;var c=t.A;return 
void 0!==c.r&&r.parentNode.removeChild(r),c.s=jn(r,t.w),r;case 8:return 
function(r,n){var t=n.s,e=function(r,n){if(r){for(var 
t=en.createDocumentFragment(),e=0;r.length>e;e++){var 
u=r[e].A;un(t,2===u.c?u.s:yn(u.z,n.u))}return 
t}}(t.y,n);r=jn(r,t.w);for(var u=t.x,a=0;u.length>a;a++){var 
c=u[a],i=c.A,o=2===i.c?i.s:yn(i.z,n.u);r.insertBefore(o,r.childNodes[c.r])}return 
e&&un(r,e),r}(r,n);case 5:return n.s(r);default:B(10)}}var 
_n=u(function(r,n,t,e){return function(r,n,t,e,u,a){var 
c=f(mr,r,kr(n?n.flags:void 0));se(c)||B(2);var 
i={},o=t(c.a),s=o.a,l=a(d,s),b=function(r,n){var t;for(var e in Gr){var 
u=Gr[e];u.a&&((t=t||{})[e]=u.a(e,n)),r[e]=Hr(u,n)}return t}(i,d);function 
d(r,n){var t=f(e,r,s);l(s=t.a,n),Qr(i,t.b,u(s))}return 
Qr(i,o.b,u(s)),b?{ports:b}:{}}(n,e,r.dg,r.ee,r.dY,function(n,t){var 
u=r.eg,a=e.node,c=function r(n){if(3===n.nodeType)return 
an(n.textContent);if(1!==n.nodeType)return an("");for(var 
t=L,e=n.attributes,u=e.length;u--;){var 
a=e[u];t=D(f(dn,a.name,a.value),t)}var 
c=n.tagName.toLowerCase(),i=L,o=n.childNodes;for(u=o.length;u--;)i=D(r(o[u]),i);return 
s(on,c,t,i)}(a);return function(r,n){n(r);var t=0;function 
e(){t=1===t?0:(Gn(e),n(r),1)}return 
function(u,a){r=u,a?(n(r),2===t&&(t=1)):(0===t&&Gn(e),t=2)}}(t,function(r){var 
t=u(r),e=function(r,n){var t=[];return 
Nn(r,n,t,0),t}(c,t);a=Pn(a,c,e,n),c=t})})}),Gn=("undefined"!=typeof 
cancelAnimationFrame&&cancelAnimationFrame,"undefined"!=typeof 
requestAnimationFrame?requestAnimationFrame:function(r){return 
setTimeout(r,1e3/60)});"undefined"!=typeof 
document&&document,"undefined"!=typeof window&&window;var 
zn=e(function(r,n,t){return Tr(function(e){function 
u(r){e(n(t.aU.a(r)))}var a=new 
XMLHttpRequest;a.addEventListener("error",function(){u(vu)}),a.addEventListener("timeout",function(){u(gu)}),a.addEventListener("load",function(){u(function(r,n){return 
f(n.status>=200&&300>n.status?du:lu,function(r){return{ef:r.responseURL,dT:r.status,dU:r.statusText,da:Hn(r.getAllResponseHeaders())}}(n),r(n.response))}(t.aU.b,a))}),wu(t.cH)&&function(r,n,t){n.upload.addEventListener("progress",function(e){n.c||Vr(f(yu,r,y(t,mu({dR:e.loaded,cv:e.total}))))}),n.addEventListener("progress",function(e){n.c||Vr(f(yu,r,y(t,pu({dJ:e.loaded,cv:e.lengthComputable?ht(e.total):$t}))))})}(r,a,t.cH.a);try{a.open(t.dm,t.ef,!0)}catch(r){return 
u(bu(t.ef))}return function(r,n){for(var 
t=n.da;t.b;t=t.b)r.setRequestHeader(t.a.a,t.a.b);r.timeout=n.d7.a||0,r.responseType=n.aU.d,r.withCredentials=n.cS}(a,t),t.cV.a&&a.setRequestHeader("Content-Type",t.cV.a),a.send(t.cV.b),function(){a.c=!0,a.abort()}})});function 
Hn(r){if(!r)return $u;for(var n=$u,t=r.split("\r\n"),e=t.length;e--;){var 
u=t[e],a=u.indexOf(": ");if(a>0){var 
c=u.substring(0,a),i=u.substring(a+2);n=s(Uu,c,function(r){return 
ht(wu(r)?i+", "+r.a:i)},n)}}return n}var 
Mn=e(function(r,n,t){return{$:0,d:r,b:n,a:t}}),Fn=t(function(r,n){return{$:0,a:r,b:n}}),Kn=a(function(r,n,t,e,u){for(var 
a=r.length,c=u.length>=n+a,i=0;c&&a>i;){var 
o=u.charCodeAt(n);c=r[i++]===u[n++]&&(10===o?(t++,e=1):(e++,55296==(63488&o)?r[i++]===u[n++]:1))}return 
q(c?n:-1,t,e)}),Zn=e(function(r,n,t){return 
t.length>n?55296==(63488&t.charCodeAt(n))?r(x(t.substr(n,2)))?n+2:-1:r(x(t[n]))?"\n"===t[n]?-2:n+1:-1:-1}),Jn=e(function(r,n,t){return 
t.charCodeAt(n)===r}),Xn=t(function(r,n){for(;n.length>r;r++){var 
t=n.charCodeAt(r);if(48>t||t>57)return r}return 
r}),Yn=e(function(r,n,t){for(var e=0;t.length>n;n++){var 
u=t.charCodeAt(n)-48;if(0>u||u>=r)break;e=r*e+u}return 
y(n,e)}),Qn=t(function(r,n){for(var t=0;n.length>r;r++){var 
e=n.charCodeAt(r);if(48>e||e>57)if(65>e||e>70){if(97>e||e>102)break;t=16*t+e-87}else 
t=16*t+e-55;else t=16*t+e-48}return 
y(r,t)}),Wn=a(function(r,n,t,e,u){for(var 
a=u.indexOf(r,n),c=0>a?u.length:a+r.length;c>n;){var 
i=u.charCodeAt(n++);10===i?(e=1,t++):(e++,55296==(63488&i)&&n++)}return 
q(a,t,e)}),rt=t(function(r,n){return Tr(function(){var 
t=setInterval(function(){Vr(n)},r);return 
function(){clearInterval(t)}})}),nt=t(function(r,n){var 
t="g";r.dp&&(t+="m"),r.cX&&(t+="i");try{return 
ht(RegExp(n,t))}catch(r){return $t}}),tt=t(function(r,n){return 
null!==n.match(r)}),et=e(function(r,n,t){for(var 
e,u=[],a=0,c=t,i=n.lastIndex,o=-1;a++<r&&(e=n.exec(c))&&o!=n.lastIndex;){for(var 
f=e.length-1,s=[];f>0;){var 
b=e[f];s[--f]=b?ht(b):$t}u.push(l(zl,e[0],e.index,a,N(s))),o=n.lastIndex}return 
n.lastIndex=i,N(u)}),ut=u(function(r,n,t,e){var u=0;return 
e.replace(n,function(n){if(u++>=r)return n;for(var 
e=arguments.length-3,a=[];e>0;){var 
c=arguments[e];a[--e]=c?ht(c):$t}return 
t(l(zl,n,arguments[arguments.length-2],u,N(a)))})}),at=1,ct=2,it=0,ot=T,ft=e(function(r,n,t){for(;;){if(-2===t.$)return 
n;var 
e=t.d,u=r,a=s(r,t.b,t.c,s(ft,r,n,t.e));r=u,n=a,t=e}}),st=function(r){return 
s(ft,e(function(r,n,t){return 
f(ot,y(r,n),t)}),L,r)},lt=function(r){return{$:1,a:r}},bt=t(function(r,n){return{$:3,a:r,b:n}}),dt=t(function(r,n){return{$:0,a:r,b:n}}),vt=t(function(r,n){return{$:1,a:r,b:n}}),pt=function(r){return{$:0,a:r}},mt=function(r){return{$:2,a:r}},gt=O,ht=function(r){return{$:0,a:r}},$t={$:1},wt=Y,yt=k,qt=Ar,xt=function(r){return 
r+""},At=t(function(r,n){return f(J,r,S(n))}),kt=t(function(r,n){return 
N(f(Z,r,n))}),Et=function(r){return f(At,"\n    
",f(kt,"\n",r))},Lt=e(function(r,n,t){for(;;){if(!t.b)return n;var 
e=t.b,u=r,a=f(r,t.a,n);r=u,n=a,t=e}}),Dt=function(r){return 
s(Lt,t(function(r,n){return 
n+1}),0,r)},Tt=C,Nt=e(function(r,n,t){for(;;){if(h(r,n)>=1)return t;var 
e=r,u=n-1,a=f(ot,n,t);r=e,n=u,t=a}}),St=t(function(r,n){return 
s(Nt,r,n,L)}),Ct=t(function(r,n){return 
s(Tt,r,f(St,0,Dt(n)-1),n)}),Vt=function(r){var n=r.charCodeAt(0);return 
55296>n||n>56319?n:1024*(n-55296)+r.charCodeAt(1)-56320+65536},Rt=function(r){var 
n=Vt(r);return n>=97&&122>=n},Ut=function(r){var n=Vt(r);return 
90>=n&&n>=65},Bt=function(r){return Rt(r)||Ut(r)},Ot=function(r){var 
n=Vt(r);return 57>=n&&n>=48},Pt=function(r){return 
Rt(r)||Ut(r)||Ot(r)},jt=function(r){return 
s(Lt,ot,L,r)},It=t(function(r,n){return"\n\n("+xt(r+1)+") 
"+Et(_t(n))}),_t=function(r){return 
f(Gt,r,L)},Gt=t(function(r,n){r:for(;;)switch(r.$){case 0:var 
t=r.a,e=r.b,u=function(){var 
r,n,e=(n=(r=t).charCodeAt(0),isNaN(n)?$t:ht(55296>n||n>56319?y(x(r[0]),r.slice(1)):y(x(r[0]+r[1]),r.slice(2))));if(1===e.$)return!1;var 
u=e.a,a=u.b;return 
Bt(u.a)&&f(wt,Pt,a)}();r=e,n=f(ot,u?"."+t:"['"+t+"']",n);continue r;case 
1:e=r.b;var a="["+xt(r.a)+"]";r=e,n=f(ot,a,n);continue r;case 2:var 
c=r.a;if(c.b){if(c.b.b){var i=(n.b?"The Json.Decode.oneOf at 
json"+f(At,"",jt(n)):"Json.Decode.oneOf")+" failed in the following 
"+xt(Dt(c))+" ways:";return 
f(At,"\n\n",f(ot,i,f(Ct,It,c)))}r=e=c.a,n=n;continue r}return"Ran into a 
Json.Decode.oneOf with no possibilities"+(n.b?" at 
json"+f(At,"",jt(n)):"!");default:var o=r.a,s=r.b;return(i=n.b?"Problem 
with the value at json"+f(At,"",jt(n))+":\n\n    ":"Problem with the given 
value:\n\n")+Et(f(qt,4,s))+"\n\n"+o}}),zt=u(function(r,n,t,e){return{$:0,a:r,b:n,c:t,d:e}}),Ht=[],Mt=_,Ft=t(function(r,n){return 
H(n)/H(r)}),Kt=Mt(f(Ft,2,32)),Zt=l(zt,0,Kt,Ht,Ht),Jt=R,Xt=t(function(r,n){return 
r(n)}),Yt=t(function(r,n){return n(r)}),Qt=g,Wt=G,re=function(r){return 
r.length},ne=t(function(r,n){return 
h(r,n)>0?r:n}),te=P,ee=U,ue=t(function(r,n){for(;;){var 
t=f(ee,32,r),e=t.b,u=f(ot,{$:0,a:t.a},n);if(!e.b)return 
jt(u);r=e,n=u}}),ae=function(r){return r.a},ce=t(function(r,n){for(;;){var 
t=Mt(n/32);if(1===t)return 
f(ee,32,r).a;r=f(ue,r,L),n=t}}),ie=t(function(r,n){if(n.h){var 
t=32*n.h,e=Wt(f(Ft,32,t-1)),u=r?jt(n.k):n.k,a=f(ce,u,n.h);return 
l(zt,re(n.i)+t,f(ne,5,e*Kt),a,n.i)}return 
l(zt,re(n.i),Kt,Ht,n.i)}),oe=a(function(r,n,t,e,u){for(;;){if(0>n)return 
f(ie,!1,{k:e,h:t/32|0,i:u});var 
a={$:1,a:s(Jt,32,n,r)};r=r,n-=32,t=t,e=f(ot,a,e),u=u}}),fe=t(function(r,n){if(r>0){var 
t=r%32;return b(oe,n,r-t-32,r,L,s(Jt,t,r-t,n))}return 
Zt}),se=function(r){return!r.$},le=sr,be=lr,de=function(r){return{$:0,a:r}},ve=function(r){switch(r.$){case 
0:return 0;case 1:return 1;case 2:return 2;default:return 
3}},pe=function(r){return r},me=Q,ge=function(r){return 
r.length},he=X,$e=t(function(r,n){return 
1>r?n:s(he,r,ge(n),n)}),we=function(r){return""===r},ye=t(function(r,n){return 
1>r?"":s(he,0,r,n)}),qe=function(r){for(var 
n=0,t=r.charCodeAt(0),e=43==t||45==t?1:0,u=e;r.length>u;++u){var 
a=r.charCodeAt(u);if(48>a||a>57)return $t;n=10*n+a-48}return 
u==e?$t:ht(45==t?-n:n)},xe=W,Ae=function(r){for(;;)r=r},ke=Dr,Ee=ke(0),Le=u(function(r,n,t,e){if(e.b){var 
u=e.a,a=e.b;if(a.b){var c=a.a,i=a.b;if(i.b){var o=i.a,b=i.b;if(b.b){var 
d=b.b;return 
f(r,u,f(r,c,f(r,o,f(r,b.a,t>500?s(Lt,r,n,jt(d)):l(Le,r,n,t+1,d)))))}return 
f(r,u,f(r,c,f(r,o,n)))}return f(r,u,f(r,c,n))}return f(r,u,n)}return 
n}),De=e(function(r,n,t){return l(Le,r,n,0,t)}),Te=t(function(r,n){return 
s(De,t(function(n,t){return 
f(ot,r(n),t)}),L,n)}),Ne=Nr,Se=t(function(r,n){return 
f(Ne,function(n){return ke(r(n))},n)}),Ce=e(function(r,n,t){return 
f(Ne,function(n){return f(Ne,function(t){return 
ke(f(r,n,t))},t)},n)}),Ve=function(r){return 
s(De,Ce(ot),ke(L),r)},Re=Mr,Ue=t(function(r,n){var t=n;return 
Rr(f(Ne,Re(r),t))});Gr.Task=zr(Ee,e(function(r,n){return 
f(Se,function(){return 0},Ve(f(Te,Ue(r),n)))}),e(function(){return 
ke(0)}),t(function(r,n){return f(Se,r,n)}));var 
Be,Oe=Kr("Task"),Pe=t(function(r,n){return 
Oe(f(Se,r,n))}),je=_n,Ie=function(r){return{$:1,a:r}},_e=function(r){return{$:6,a:r}},Ge=function(r){return{$:0,a:r}},ze=e(function(r,n,t){return 
r(n(t))}),He=Sr,Me=t(function(r,n){return 
Oe(f(He,f(ze,f(ze,ke,r),lt),f(Ne,f(ze,f(ze,ke,r),pt),n)))}),Fe=Zr,Ke=function(r){return{$:1,a:r}},Ze=t(function(r,n){return{$:0,a:r,b:n}}),Je=fr,Xe=ir,Ye=tr,Qe=function(r){return{$:11,g:r}},We=t(function(r,n){return 
y(r,n)}),ru=cr,nu=f(Je,function(r){var 
n=r.a,t=r.b;return"M_LIMIT_EXCEEDED"===n?Qe(N([f(le,Ke,f(Xe,"retry_after_ms",Ye)),de(f(Ze,n,t))])):de(f(Ze,n,t))},s(be,We,f(Xe,"errcode",ru),f(Xe,"error",ru))),tu=pr,eu=t(function(r,n){return 
n.$?r:n.a}),uu=t(function(r,n){switch(n.$){case 0:return lt(f(Ze,"Bad 
url",n.a));case 1:case 2:return lt(Ke(30));case 3:return 
lt(f(eu,f(Ze,"Could not decode error",t=n.b),f(tu,nu,t)));default:var 
t,e=f(tu,r,t=n.b);return 1===e.$?lt(f(Ze,"Could not decode 
response",t)):pt(e.a)}}),au=t(function(r,n){return{$:0,a:r,b:n}}),cu=t(function(r,n){return 
n.$?$t:ht(r(n.a))}),iu=t(function(r,n){return{$:0,a:r,b:n}}),ou=function(r){return{$:1,a:r}},fu=_r,su=t(function(r,n){return 
n.$?f(Ne,function(){return 
f(He,su(r),r)},fu(n.a)):ou(f(iu,n.a,n.b))}),lu=t(function(r,n){return{$:3,a:r,b:n}}),bu=function(r){return{$:0,a:r}},du=t(function(r,n){return{$:4,a:r,b:n}}),vu={$:2},pu=function(r){return{$:1,a:r}},mu=function(r){return{$:0,a:r}},gu={$:1},hu={$:-2},$u=hu,wu=function(r){return!r.$},yu=Fr,qu=$,xu=t(function(r,n){r:for(;;){if(-2===n.$)return 
$t;var t=n.c,e=n.d,u=n.e;switch(f(qu,r,n.b)){case 0:r=r,n=e;continue 
r;case 1:return ht(t);default:r=r,n=u;continue 
r}}}),Au=a(function(r,n,t,e,u){return{$:-1,a:r,b:n,c:t,d:e,e:u}}),ku=a(function(r,n,t,e,u){if(-1!==u.$||u.a){if(-1!==e.$||e.a||-1!==e.d.$||e.d.a)return 
b(Au,r,n,t,e,u);var a=e.d;return 
c=e.e,b(Au,0,e.b,e.c,b(Au,1,a.b,a.c,a.d,a.e),b(Au,1,n,t,c,u))}var 
c,i=u.b,o=u.c,f=u.d,s=u.e;return-1!==e.$||e.a?b(Au,r,i,o,b(Au,0,n,t,e,f),s):b(Au,0,n,t,b(Au,1,e.b,e.c,e.d,c=e.e),b(Au,1,i,o,f,s))}),Eu=e(function(r,n,t){if(-2===t.$)return 
b(Au,0,r,n,hu,hu);var e=t.a,u=t.b,a=t.c,c=t.d,i=t.e;switch(f(qu,r,u)){case 
0:return b(ku,e,u,a,s(Eu,r,n,c),i);case 1:return 
b(Au,e,u,n,c,i);default:return 
b(ku,e,u,a,c,s(Eu,r,n,i))}}),Lu=e(function(r,n,t){var 
e=s(Eu,r,n,t);return-1!==e.$||e.a?e:b(Au,1,e.b,e.c,e.d,e.e)}),Du=function(r){if(-1===r.$&&-1===r.d.$&&-1===r.e.$){if(-1!==r.e.d.$||r.e.d.a){var 
n=r.d,t=r.e;return 
c=t.b,i=t.c,e=t.d,s=t.e,b(Au,1,r.b,r.c,b(Au,0,n.b,n.c,n.d,n.e),b(Au,0,c,i,e,s))}var 
e,u=r.d,a=r.e,c=a.b,i=a.c,o=(e=a.d).d,f=e.e,s=a.e;return 
b(Au,0,e.b,e.c,b(Au,1,r.b,r.c,b(Au,0,u.b,u.c,u.d,u.e),o),b(Au,1,c,i,f,s))}return 
r},Tu=function(r){if(-1===r.$&&-1===r.d.$&&-1===r.e.$){if(-1!==r.d.d.$||r.d.d.a){var 
n=r.d,t=r.e;return 
f=t.b,s=t.c,l=t.d,d=t.e,b(Au,1,e=r.b,u=r.c,b(Au,0,n.b,n.c,n.d,i=n.e),b(Au,0,f,s,l,d))}var 
e=r.b,u=r.c,a=r.d,c=a.d,i=a.e,o=r.e,f=o.b,s=o.c,l=o.d,d=o.e;return 
b(Au,0,a.b,a.c,b(Au,1,c.b,c.c,c.d,c.e),b(Au,1,e,u,i,b(Au,0,f,s,l,d)))}return 
r},Nu=i(function(r,n,t,e,u,a,c){if(-1!==a.$||a.a){r:for(;;){if(-1===c.$&&1===c.a){if(-1===c.d.$){if(1===c.d.a)return 
Tu(n);break r}return Tu(n)}break r}return n}return 
b(Au,t,a.b,a.c,a.d,b(Au,0,e,u,a.e,c))}),Su=function(r){if(-1===r.$&&-1===r.d.$){var 
n=r.a,t=r.b,e=r.c,u=r.d,a=u.d,c=r.e;if(1===u.a){if(-1!==a.$||a.a){var 
i=Du(r);if(-1===i.$){var o=i.e;return b(ku,i.a,i.b,i.c,Su(i.d),o)}return 
hu}return b(Au,n,t,e,Su(u),c)}return b(Au,n,t,e,Su(u),c)}return 
hu},Cu=t(function(r,n){if(-2===n.$)return hu;var 
t=n.a,e=n.b,u=n.c,a=n.d,c=n.e;if(0>h(r,e)){if(-1===a.$&&1===a.a){var 
i=a.d;if(-1!==i.$||i.a){var o=Du(n);if(-1===o.$){var s=o.e;return 
b(ku,o.a,o.b,o.c,f(Cu,r,o.d),s)}return hu}return 
b(Au,t,e,u,f(Cu,r,a),c)}return b(Au,t,e,u,f(Cu,r,a),c)}return 
f(Vu,r,v(Nu,r,n,t,e,u,a,c))}),Vu=t(function(r,n){if(-1===n.$){var 
t=n.a,e=n.b,u=n.c,a=n.d,c=n.e;if(p(r,e)){var 
i=function(r){for(;;){if(-1!==r.$||-1!==r.d.$)return 
r;r=r.d}}(c);return-1===i.$?b(ku,t,i.b,i.c,a,Su(c)):hu}return 
b(ku,t,e,u,a,f(Cu,r,c))}return hu}),Ru=t(function(r,n){var 
t=f(Cu,r,n);return-1!==t.$||t.a?t:b(Au,1,t.b,t.c,t.d,t.e)}),Uu=e(function(r,n,t){var 
e=n(f(xu,r,t));return 
e.$?f(Ru,r,t):s(Lu,r,e.a,t)}),Bu=f(Mn,"",pe),Ou=function(r){return 
r.$?ou(r.a):ke(r.a)},Pu=function(r){return 
s(zn,0,Ou,{cS:!1,cV:r.cV,aU:r.dL,da:r.da,dm:r.dm,d7:r.d7,cH:$t,ef:r.ef})},ju=t(function(r,n){return 
n.$?r:n.a}),Iu=function(r){var 
n=r.dm,t=r.ef,e=r.dM,u=Pu({cV:r.cV,da:f(ju,L,f(cu,function(r){return 
N([f(au,"Authorization","Bearer 
"+r)])},r.U)),dm:n,dL:Bu(uu(e)),d7:$t,ef:t});return 
f(He,su(u),u)},_u=function(r){return 
r.a+"="+r.b},Gu=e(function(r,n,t){return 
r+"/"+(f(At,"/",n)+function(r){return 
r.b?"?"+f(At,"&",f(Te,_u,r)):""}(t))}),zu=function(r){return 
encodeURIComponent(r)},Hu=u(function(r,n,t,e){return 
s(Gu,n,f(Te,zu,E(r,t)),e)}),Mu=Hu(N(["_matrix","client","r0"])),Fu=t(function(r,n){var 
t=r,e=n.dm,u=n.dG,a=n.dF,c=n.cV,i=n.dM;return 
Iu({U:ht(t.U),cV:c,dm:e,dM:i,ef:s(Mu,t.bT,u,a)})}),Ku=t(function(r,n){return{$:1,a:r,b:n}}),Zu=function(r){return{$:1,a:r}},Ju=f(Je,function(r){switch(r){case"invite":return 
de(0);case"join":return de(1);case"leave":return de(2);case"ban":return 
de(3);case"knock":return de(4);default:return Zu("Invalid membership 
field: "+r)}},ru),Xu=br,Yu=function(r){return 
Qe(N([f(le,ht,r),de($t)]))},Qu=l(Xu,e(function(r,n,t){return{bw:t,bI:n,b3:r}}),f(Xe,"membership",Ju),Yu(f(Xe,"displayname",ru)),Yu(f(Xe,"avatar_url",ru))),Wu=function(r){return{$:6,a:r}},ra=function(r){return{$:1,a:r}},na=function(r){return{$:4,a:r}},ta=function(r){return{$:3,a:r}},ea=function(r){return{$:2,a:r}},ua=function(r){return{$:0,a:r}},aa=function(r){return{$:5,a:r}},ca=s(be,t(function(r,n){return{cV:r,ef:n}}),f(Xe,"body",ru),f(Xe,"url",ru)),ia=s(be,t(function(r,n){return{cV:r,ef:n}}),f(Xe,"body",ru),f(Xe,"url",ru)),oa=f(le,function(r){return{$:0,a:r}},f(Xe,"body",ru)),fa=e(function(r,n,t){return{$:1,a:r,b:n,c:t}}),sa=t(function(r,n){return{$:1,a:r,b:n}}),la=e(function(r,n,t){return{$:0,a:r,b:n,c:t}}),ba=t(function(r,n){var 
t=n;return function(n){var e=t(n);if(1===e.$)return f(sa,e.a,e.b);var 
u=e.a,a=e.c,c=r(e.b)(a);if(1===c.$){var i=c.a;return f(sa,u||i,c.b)}return 
i=c.a,s(la,u||i,c.b,c.c)}}),da=ba,va=function(r){var n=r;return 
function(r){var t=n(r);return 
1===t.$?f(sa,!1,t.b):s(la,!1,t.b,t.c)}},pa=va,ma={$:11},ga=t(function(r,n){return{$:1,a:r,b:n}}),ha=u(function(r,n,t,e){return{bB:n,c0:e,dH:t,dP:r}}),$a={$:0},wa=t(function(r,n){return 
f(ga,$a,l(ha,r.dP,r.bB,n,r.d))}),ya=Zn,qa=function(r){return-r},xa=t(function(r,n){return 
function(t){var e=s(ya,r,t.b,t.cy);return 
p(e,-1)?f(sa,!1,f(wa,t,n)):p(e,-2)?s(la,!0,0,{bB:1,d:t.d,e:t.e,b:t.b+1,dP:t.dP+1,cy:t.cy}):s(la,!0,0,{bB:t.bB+1,d:t.d,e:t.e,b:e,dP:t.dP,cy:t.cy})}}),Aa=function(r){return 
f(xa,r,ma)},ka=a(function(r,n,t,e,u){for(;;){var 
a=s(ya,r,n,u.cy);if(p(a,-1))return 
s(la,0>h(u.b,n),0,{bB:e,d:u.d,e:u.e,b:n,dP:t,cy:u.cy});p(a,-2)?(r=r,n+=1,t+=1,e=1,u=u):(r=r,n=a,t=t,e+=1,u=u)}}),Ea=function(r){return 
function(n){return b(ka,r,n.b,n.dP,n.bB,n)}},La=Ea,Da=t(function(r){return 
r}),Ta=e(function(r,n,t){var e=n,u=t;return function(n){var 
t=e(n);if(1===t.$)return f(sa,t.a,t.b);var 
a=t.a,c=t.b,i=u(t.c);if(1===i.$){var o=i.a;return f(sa,a||o,i.b)}o=i.a;var 
l=i.c;return s(la,a||o,f(r,c,i.b),l)}}),Na=t(function(r,n){return 
s(Ta,Da,r,n)}),Sa=Na,Ca=function(r){return 
f(Sa,Aa(r),La(r))},Va=t(function(r,n){var t=n;return function(n){var 
e=t(n);if(1===e.$)return f(sa,e.a,e.b);var u=e.b,a=e.c;return 
s(la,e.a,f(r,s(he,n.b,a.b,n.cy),u),a)}}),Ra=function(r){return 
f(Va,Da,r)},Ua=Ra,Ba=function(r){return" 
"===r||"\t"===r||"\n"===r||"\r"===r||"\f"===r||" 
"===r},Oa=function(r){return!r},Pa=function(r){return{$:12,a:r}},ja=function(r){return 
function(n){return f(sa,!1,f(wa,n,r))}},Ia=function(r){return 
ja(Pa(r))},_a=function(r){return function(n){return 
s(la,!1,r,n)}},Ga=_a,za=function(r){return 
r.toLowerCase()},Ha=function(r){var n=f(da,function(n){return 
p(za(n),r)?Ga(0):Ia("closing tag does not match opening tag: 
"+r)},Ua(Ca(function(r){return!Ba(r)&&">"!==r})));return 
f(Sa,f(Sa,f(Sa,f(Sa,Aa(Qt("<")),Aa(Qt("/"))),n),La(Ba)),Aa(Qt(">")))},Ma=function(r){return{$:2,a:r}},Fa=Wn,Ka=u(function(r,n,t,e){return 
f(ga,$a,l(ha,r,n,t,e))}),Za=function(r){var n=r.a,t=r.b;return 
function(r){var e=b(Fa,n,r.b,r.dP,r.bB,r.cy),u=e.a,a=e.b,c=e.c;return 
p(u,-1)?f(sa,!1,l(Ka,a,c,t,r.d)):s(la,0>h(r.b,u),0,{bB:c,d:r.d,e:r.e,b:u,dP:a,cy:r.cy})}},Ja=function(r){return{$:0,a:r}},Xa=t(function(r,n){return{$:0,a:r,b:n}}),Ya=function(r){return 
f(Xa,r,Ja(r))},Qa=t(function(r,n){return 
s(Ta,Xt,r,n)}),Wa=Qa,rc=Kn,nc=function(r){var n=r.a,t=r.b,e=!we(n);return 
function(r){var u=b(rc,n,r.b,r.dP,r.bB,r.cy),a=u.a,c=u.b,i=u.c;return 
p(a,-1)?f(sa,!1,f(wa,r,t)):s(la,e,0,{bB:i,d:r.d,e:r.e,b:a,dP:c,cy:r.cy})}},tc=function(r){return 
nc(Ya(r))},ec=f(Wa,f(Sa,f(Sa,Ga(pe),tc("<!")),tc("--")),f(Sa,Ua(Za(Ya("--\x3e"))),tc("--\x3e"))),uc=t(function(r,n){var 
t=n;return function(n){var e=t(n);if(e.$)return f(sa,e.a,e.b);var 
u=e.c;return 
s(la,e.a,r(e.b),u)}}),ac=uc,cc=f(ac,Ma,ec),ic=t(function(r,n){for(;;){if(!n.b)return!1;var 
t=n.b;if(r(n.a))return!0;r=r,n=t}}),oc=t(function(r,n){return 
f(ic,function(n){return 
p(n,r)},n)}),fc=N(["area","base","br","col","embed","hr","img","input","link","meta","param","source","track","wbr"]),sc=function(r){return 
f(oc,r,fc)},lc=function(r){return{$:1,a:r}},bc=function(r){return{$:0,a:r}},dc=u(function(r,n,t,e){for(;;){var 
u=t(n)(e);if(u.$)return a=u.a,f(sa,r||a,u.b);var 
a=u.a,c=u.b,i=u.c;if(c.$)return 
s(la,r||a,c.a,i);r=r||a,n=c.a,t=t,e=i}}),vc=t(function(r,n){return 
function(t){return 
l(dc,!1,r,n,t)}}),pc=function(r){return{$:1,a:r}},mc=function(r){return{$:0,a:r}},gc=function(r){return 
r.$?pc(r.a):mc(r.a)},hc=t(function(r,n){return f(vc,r,function(r){return 
f(ac,gc,n(r))})}),$c=t(function(r,n){return{$:2,a:r,b:n}}),wc=e(function(r,n,t){r:for(;;){if(t.b){var 
e=t.b,u=(0,t.a)(r);if(u.$){var a;if((a=u).a)return 
a;r=r,n=f($c,n,a.b),t=e;continue r}return u}return 
f(sa,!1,n)}}),yc=function(r){return function(n){return 
s(wc,n,$a,r)}},qc=yc,xc=function(r){return f(hc,L,function(n){return 
qc(N([f(ac,function(r){return 
bc(f(ot,r,n))},r),Ga(lc(jt(n)))]))})},Ac=f(ac,za,Ua(Ca(function(r){return!Ba(r)&&'"'!==r&&"'"!==r&&">"!==r&&"/"!==r&&"="!==r}))),kc=Aa(Qt(";")),Ec=function(r){return 
s(Lt,t(function(r,n){return 
s(Lu,r.a,r.b,n)}),$u,r)},Lc=Ec(N([y("Aacute","Ã"),y("aacute","Ã¡"),y("Abreve","Ä‚"),y("abreve","Äƒ"),y("ac","âˆ¾"),y("acd","âˆ¿"),y("acE","âˆ¾Ì³"),y("Acirc","Ã‚"),y("acirc","Ã¢"),y("acute","Â´"),y("Acy","Ð"),y("acy","Ð°"),y("AElig","Ã†"),y("aelig","Ã¦"),y("af","â¡"),y("Afr","ð”„"),y("afr","ð”ž"),y("Agrave","Ã€"),y("agrave","Ã 
"),y("alefsym","â„µ"),y("aleph","â„µ"),y("Alpha","Î‘"),y("alpha","Î±"),y("Amacr","Ä€"),y("amacr","Ä"),y("amalg","â¨¿"),y("amp","&"),y("AMP","&"),y("andand","â©•"),y("And","â©“"),y("and","âˆ§"),y("andd","â©œ"),y("andslope","â©˜"),y("andv","â©š"),y("ang","âˆ 
"),y("ange","â¦¤"),y("angle","âˆ 
"),y("angmsdaa","â¦¨"),y("angmsdab","â¦©"),y("angmsdac","â¦ª"),y("angmsdad","â¦«"),y("angmsdae","â¦¬"),y("angmsdaf","â¦­"),y("angmsdag","â¦®"),y("angmsdah","â¦¯"),y("angmsd","âˆ¡"),y("angrt","âˆŸ"),y("angrtvb","âŠ¾"),y("angrtvbd","â¦"),y("angsph","âˆ¢"),y("angst","Ã…"),y("angzarr","â¼"),y("Aogon","Ä„"),y("aogon","Ä…"),y("Aopf","ð”¸"),y("aopf","ð•’"),y("apacir","â©¯"),y("ap","â‰ˆ"),y("apE","â©°"),y("ape","â‰Š"),y("apid","â‰‹"),y("apos","'"),y("ApplyFunction","â¡"),y("approx","â‰ˆ"),y("approxeq","â‰Š"),y("Aring","Ã…"),y("aring","Ã¥"),y("Ascr","ð’œ"),y("ascr","ð’¶"),y("Assign","â‰”"),y("ast","*"),y("asymp","â‰ˆ"),y("asympeq","â‰"),y("Atilde","Ãƒ"),y("atilde","Ã£"),y("Auml","Ã„"),y("auml","Ã¤"),y("awconint","âˆ³"),y("awint","â¨‘"),y("backcong","â‰Œ"),y("backepsilon","Ï¶"),y("backprime","â€µ"),y("backsim","âˆ½"),y("backsimeq","â‹"),y("Backslash","âˆ–"),y("Barv","â«§"),y("barvee","âŠ½"),y("barwed","âŒ…"),y("Barwed","âŒ†"),y("barwedge","âŒ…"),y("bbrk","âŽµ"),y("bbrktbrk","âŽ¶"),y("bcong","â‰Œ"),y("Bcy","Ð‘"),y("bcy","Ð±"),y("bdquo","â€ž"),y("becaus","âˆµ"),y("because","âˆµ"),y("Because","âˆµ"),y("bemptyv","â¦°"),y("bepsi","Ï¶"),y("bernou","â„¬"),y("Bernoullis","â„¬"),y("Beta","Î’"),y("beta","Î²"),y("beth","â„¶"),y("between","â‰¬"),y("Bfr","ð”…"),y("bfr","ð”Ÿ"),y("bigcap","â‹‚"),y("bigcirc","â—¯"),y("bigcup","â‹ƒ"),y("bigodot","â¨€"),y("bigoplus","â¨"),y("bigotimes","â¨‚"),y("bigsqcup","â¨†"),y("bigstar","â˜…"),y("bigtriangledown","â–½"),y("bigtriangleup","â–³"),y("biguplus","â¨„"),y("bigvee","â‹"),y("bigwedge","â‹€"),y("bkarow","â¤"),y("blacklozenge","â§«"),y("blacksquare","â–ª"),y("blacktriangle","â–´"),y("blacktriangledown","â–¾"),y("blacktriangleleft","â—‚"),y("blacktriangleright","â–¸"),y("blank","â£"),y("blk12","â–’"),y("blk14","â–‘"),y("blk34","â–“"),y("block","â–ˆ"),y("bne","=âƒ¥"),y("bnequiv","â‰¡âƒ¥"),y("bNot","â«­"),y("bnot","âŒ"),y("Bopf","ð”¹"),y("bopf","ð•“"),y("bot","âŠ¥"),y("bottom","âŠ¥"),y("bowtie","â‹ˆ"),y("boxbox","â§‰"),y("boxdl","â”"),y("boxdL","â••"),y("boxDl","â•–"),y("boxDL","â•—"),y("boxdr","â”Œ"),y("boxdR","â•’"),y("boxDr","â•“"),y("boxDR","â•”"),y("boxh","â”€"),y("boxH","â•"),y("boxhd","â”¬"),y("boxHd","â•¤"),y("boxhD","â•¥"),y("boxHD","â•¦"),y("boxhu","â”´"),y("boxHu","â•§"),y("boxhU","â•¨"),y("boxHU","â•©"),y("boxminus","âŠŸ"),y("boxplus","âŠž"),y("boxtimes","âŠ 
"),y("boxul","â”˜"),y("boxuL","â•›"),y("boxUl","â•œ"),y("boxUL","â•"),y("boxur","â””"),y("boxuR","â•˜"),y("boxUr","â•™"),y("boxUR","â•š"),y("boxv","â”‚"),y("boxV","â•‘"),y("boxvh","â”¼"),y("boxvH","â•ª"),y("boxVh","â•«"),y("boxVH","â•¬"),y("boxvl","â”¤"),y("boxvL","â•¡"),y("boxVl","â•¢"),y("boxVL","â•£"),y("boxvr","â”œ"),y("boxvR","â•ž"),y("boxVr","â•Ÿ"),y("boxVR","â• 
"),y("bprime","â€µ"),y("breve","Ë˜"),y("Breve","Ë˜"),y("brvbar","Â¦"),y("bscr","ð’·"),y("Bscr","â„¬"),y("bsemi","â"),y("bsim","âˆ½"),y("bsime","â‹"),y("bsolb","â§…"),y("bsol","\\"),y("bsolhsub","âŸˆ"),y("bull","â€¢"),y("bullet","â€¢"),y("bump","â‰Ž"),y("bumpE","âª®"),y("bumpe","â‰"),y("Bumpeq","â‰Ž"),y("bumpeq","â‰"),y("Cacute","Ä†"),y("cacute","Ä‡"),y("capand","â©„"),y("capbrcup","â©‰"),y("capcap","â©‹"),y("cap","âˆ©"),y("Cap","â‹’"),y("capcup","â©‡"),y("capdot","â©€"),y("CapitalDifferentialD","â……"),y("caps","âˆ©ï¸€"),y("caret","â"),y("caron","Ë‡"),y("Cayleys","â„­"),y("ccaps","â©"),y("Ccaron","ÄŒ"),y("ccaron","Ä"),y("Ccedil","Ã‡"),y("ccedil","Ã§"),y("Ccirc","Äˆ"),y("ccirc","Ä‰"),y("Cconint","âˆ°"),y("ccups","â©Œ"),y("ccupssm","â©"),y("Cdot","ÄŠ"),y("cdot","Ä‹"),y("cedil","Â¸"),y("Cedilla","Â¸"),y("cemptyv","â¦²"),y("cent","Â¢"),y("centerdot","Â·"),y("CenterDot","Â·"),y("cfr","ð” 
"),y("Cfr","â„­"),y("CHcy","Ð§"),y("chcy","Ñ‡"),y("check","âœ“"),y("checkmark","âœ“"),y("Chi","Î§"),y("chi","Ï‡"),y("circ","Ë†"),y("circeq","â‰—"),y("circlearrowleft","â†º"),y("circlearrowright","â†»"),y("circledast","âŠ›"),y("circledcirc","âŠš"),y("circleddash","âŠ"),y("CircleDot","âŠ™"),y("circledR","Â®"),y("circledS","â“ˆ"),y("CircleMinus","âŠ–"),y("CirclePlus","âŠ•"),y("CircleTimes","âŠ—"),y("cir","â—‹"),y("cirE","â§ƒ"),y("cire","â‰—"),y("cirfnint","â¨"),y("cirmid","â«¯"),y("cirscir","â§‚"),y("ClockwiseContourIntegral","âˆ²"),y("CloseCurlyDoubleQuote","â€"),y("CloseCurlyQuote","â€™"),y("clubs","â™£"),y("clubsuit","â™£"),y("colon",":"),y("Colon","âˆ·"),y("Colone","â©´"),y("colone","â‰”"),y("coloneq","â‰”"),y("comma",","),y("commat","@"),y("comp","âˆ"),y("compfn","âˆ˜"),y("complement","âˆ"),y("complexes","â„‚"),y("cong","â‰…"),y("congdot","â©­"),y("Congruent","â‰¡"),y("conint","âˆ®"),y("Conint","âˆ¯"),y("ContourIntegral","âˆ®"),y("copf","ð•”"),y("Copf","â„‚"),y("coprod","âˆ"),y("Coproduct","âˆ"),y("copy","Â©"),y("COPY","Â©"),y("copysr","â„—"),y("CounterClockwiseContourIntegral","âˆ³"),y("crarr","â†µ"),y("cross","âœ—"),y("Cross","â¨¯"),y("Cscr","ð’ž"),y("cscr","ð’¸"),y("csub","â«"),y("csube","â«‘"),y("csup","â«"),y("csupe","â«’"),y("ctdot","â‹¯"),y("cudarrl","â¤¸"),y("cudarrr","â¤µ"),y("cuepr","â‹ž"),y("cuesc","â‹Ÿ"),y("cularr","â†¶"),y("cularrp","â¤½"),y("cupbrcap","â©ˆ"),y("cupcap","â©†"),y("CupCap","â‰"),y("cup","âˆª"),y("Cup","â‹“"),y("cupcup","â©Š"),y("cupdot","âŠ"),y("cupor","â©…"),y("cups","âˆªï¸€"),y("curarr","â†·"),y("curarrm","â¤¼"),y("curlyeqprec","â‹ž"),y("curlyeqsucc","â‹Ÿ"),y("curlyvee","â‹Ž"),y("curlywedge","â‹"),y("curren","Â¤"),y("curvearrowleft","â†¶"),y("curvearrowright","â†·"),y("cuvee","â‹Ž"),y("cuwed","â‹"),y("cwconint","âˆ²"),y("cwint","âˆ±"),y("cylcty","âŒ­"),y("dagger","â€ 
"),y("Dagger","â€¡"),y("daleth","â„¸"),y("darr","â†“"),y("Darr","â†¡"),y("dArr","â‡“"),y("dash","â€"),y("Dashv","â«¤"),y("dashv","âŠ£"),y("dbkarow","â¤"),y("dblac","Ë"),y("Dcaron","ÄŽ"),y("dcaron","Ä"),y("Dcy","Ð”"),y("dcy","Ð´"),y("ddagger","â€¡"),y("ddarr","â‡Š"),y("DD","â……"),y("dd","â…†"),y("DDotrahd","â¤‘"),y("ddotseq","â©·"),y("deg","Â°"),y("Del","âˆ‡"),y("Delta","Î”"),y("delta","Î´"),y("demptyv","â¦±"),y("dfisht","â¥¿"),y("Dfr","ð”‡"),y("dfr","ð”¡"),y("dHar","â¥¥"),y("dharl","â‡ƒ"),y("dharr","â‡‚"),y("DiacriticalAcute","Â´"),y("DiacriticalDot","Ë™"),y("DiacriticalDoubleAcute","Ë"),y("DiacriticalGrave","`"),y("DiacriticalTilde","Ëœ"),y("diam","â‹„"),y("diamond","â‹„"),y("Diamond","â‹„"),y("diamondsuit","â™¦"),y("diams","â™¦"),y("die","Â¨"),y("DifferentialD","â…†"),y("digamma","Ï"),y("disin","â‹²"),y("div","Ã·"),y("divide","Ã·"),y("divideontimes","â‹‡"),y("divonx","â‹‡"),y("DJcy","Ð‚"),y("djcy","Ñ’"),y("dlcorn","âŒž"),y("dlcrop","âŒ"),y("dollar","$"),y("Dopf","ð”»"),y("dopf","ð••"),y("Dot","Â¨"),y("dot","Ë™"),y("DotDot","âƒœ"),y("doteq","â‰"),y("doteqdot","â‰‘"),y("DotEqual","â‰"),y("dotminus","âˆ¸"),y("dotplus","âˆ”"),y("dotsquare","âŠ¡"),y("doublebarwedge","âŒ†"),y("DoubleContourIntegral","âˆ¯"),y("DoubleDot","Â¨"),y("DoubleDownArrow","â‡“"),y("DoubleLeftArrow","â‡"),y("DoubleLeftRightArrow","â‡”"),y("DoubleLeftTee","â«¤"),y("DoubleLongLeftArrow","âŸ¸"),y("DoubleLongLeftRightArrow","âŸº"),y("DoubleLongRightArrow","âŸ¹"),y("DoubleRightArrow","â‡’"),y("DoubleRightTee","âŠ¨"),y("DoubleUpArrow","â‡‘"),y("DoubleUpDownArrow","â‡•"),y("DoubleVerticalBar","âˆ¥"),y("DownArrowBar","â¤“"),y("downarrow","â†“"),y("DownArrow","â†“"),y("Downarrow","â‡“"),y("DownArrowUpArrow","â‡µ"),y("DownBreve","Ì‘"),y("downdownarrows","â‡Š"),y("downharpoonleft","â‡ƒ"),y("downharpoonright","â‡‚"),y("DownLeftRightVector","â¥"),y("DownLeftTeeVector","â¥ž"),y("DownLeftVectorBar","â¥–"),y("DownLeftVector","â†½"),y("DownRightTeeVector","â¥Ÿ"),y("DownRightVectorBar","â¥—"),y("DownRightVector","â‡"),y("DownTeeArrow","â†§"),y("DownTee","âŠ¤"),y("drbkarow","â¤"),y("drcorn","âŒŸ"),y("drcrop","âŒŒ"),y("Dscr","ð’Ÿ"),y("dscr","ð’¹"),y("DScy","Ð…"),y("dscy","Ñ•"),y("dsol","â§¶"),y("Dstrok","Ä"),y("dstrok","Ä‘"),y("dtdot","â‹±"),y("dtri","â–¿"),y("dtrif","â–¾"),y("duarr","â‡µ"),y("duhar","â¥¯"),y("dwangle","â¦¦"),y("DZcy","Ð"),y("dzcy","ÑŸ"),y("dzigrarr","âŸ¿"),y("Eacute","Ã‰"),y("eacute","Ã©"),y("easter","â©®"),y("Ecaron","Äš"),y("ecaron","Ä›"),y("Ecirc","ÃŠ"),y("ecirc","Ãª"),y("ecir","â‰–"),y("ecolon","â‰•"),y("Ecy","Ð­"),y("ecy","Ñ"),y("eDDot","â©·"),y("Edot","Ä–"),y("edot","Ä—"),y("eDot","â‰‘"),y("ee","â…‡"),y("efDot","â‰’"),y("Efr","ð”ˆ"),y("efr","ð”¢"),y("eg","âªš"),y("Egrave","Ãˆ"),y("egrave","Ã¨"),y("egs","âª–"),y("egsdot","âª˜"),y("el","âª™"),y("Element","âˆˆ"),y("elinters","â§"),y("ell","â„“"),y("els","âª•"),y("elsdot","âª—"),y("Emacr","Ä’"),y("emacr","Ä“"),y("empty","âˆ…"),y("emptyset","âˆ…"),y("EmptySmallSquare","â—»"),y("emptyv","âˆ…"),y("EmptyVerySmallSquare","â–«"),y("emsp13","â€„"),y("emsp14","â€…"),y("emsp","â€ƒ"),y("ENG","ÅŠ"),y("eng","Å‹"),y("ensp","â€‚"),y("Eogon","Ä˜"),y("eogon","Ä™"),y("Eopf","ð”¼"),y("eopf","ð•–"),y("epar","â‹•"),y("eparsl","â§£"),y("eplus","â©±"),y("epsi","Îµ"),y("Epsilon","Î•"),y("epsilon","Îµ"),y("epsiv","Ïµ"),y("eqcirc","â‰–"),y("eqcolon","â‰•"),y("eqsim","â‰‚"),y("eqslantgtr","âª–"),y("eqslantless","âª•"),y("Equal","â©µ"),y("equals","="),y("EqualTilde","â‰‚"),y("equest","â‰Ÿ"),y("Equilibrium","â‡Œ"),y("equiv","â‰¡"),y("equivDD","â©¸"),y("eqvparsl","â§¥"),y("erarr","â¥±"),y("erDot","â‰“"),y("escr","â„¯"),y("Escr","â„°"),y("esdot","â‰"),y("Esim","â©³"),y("esim","â‰‚"),y("Eta","Î—"),y("eta","Î·"),y("ETH","Ã"),y("eth","Ã°"),y("Euml","Ã‹"),y("euml","Ã«"),y("euro","â‚¬"),y("excl","!"),y("exist","âˆƒ"),y("Exists","âˆƒ"),y("expectation","â„°"),y("exponentiale","â…‡"),y("ExponentialE","â…‡"),y("fallingdotseq","â‰’"),y("Fcy","Ð¤"),y("fcy","Ñ„"),y("female","â™€"),y("ffilig","ï¬ƒ"),y("fflig","ï¬€"),y("ffllig","ï¬„"),y("Ffr","ð”‰"),y("ffr","ð”£"),y("filig","ï¬"),y("FilledSmallSquare","â—¼"),y("FilledVerySmallSquare","â–ª"),y("fjlig","fj"),y("flat","â™­"),y("fllig","ï¬‚"),y("fltns","â–±"),y("fnof","Æ’"),y("Fopf","ð”½"),y("fopf","ð•—"),y("forall","âˆ€"),y("ForAll","âˆ€"),y("fork","â‹”"),y("forkv","â«™"),y("Fouriertrf","â„±"),y("fpartint","â¨"),y("frac12","Â½"),y("frac13","â…“"),y("frac14","Â¼"),y("frac15","â…•"),y("frac16","â…™"),y("frac18","â…›"),y("frac23","â…”"),y("frac25","â…–"),y("frac34","Â¾"),y("frac35","â…—"),y("frac38","â…œ"),y("frac45","â…˜"),y("frac56","â…š"),y("frac58","â…"),y("frac78","â…ž"),y("frasl","â„"),y("frown","âŒ¢"),y("fscr","ð’»"),y("Fscr","â„±"),y("gacute","Çµ"),y("Gamma","Î“"),y("gamma","Î³"),y("Gammad","Ïœ"),y("gammad","Ï"),y("gap","âª†"),y("Gbreve","Äž"),y("gbreve","ÄŸ"),y("Gcedil","Ä¢"),y("Gcirc","Äœ"),y("gcirc","Ä"),y("Gcy","Ð“"),y("gcy","Ð³"),y("Gdot","Ä 
"),y("gdot","Ä¡"),y("ge","â‰¥"),y("gE","â‰§"),y("gEl","âªŒ"),y("gel","â‹›"),y("geq","â‰¥"),y("geqq","â‰§"),y("geqslant","â©¾"),y("gescc","âª©"),y("ges","â©¾"),y("gesdot","âª€"),y("gesdoto","âª‚"),y("gesdotol","âª„"),y("gesl","â‹›ï¸€"),y("gesles","âª”"),y("Gfr","ð”Š"),y("gfr","ð”¤"),y("gg","â‰«"),y("Gg","â‹™"),y("ggg","â‹™"),y("gimel","â„·"),y("GJcy","Ðƒ"),y("gjcy","Ñ“"),y("gla","âª¥"),y("gl","â‰·"),y("glE","âª’"),y("glj","âª¤"),y("gnap","âªŠ"),y("gnapprox","âªŠ"),y("gne","âªˆ"),y("gnE","â‰©"),y("gneq","âªˆ"),y("gneqq","â‰©"),y("gnsim","â‹§"),y("Gopf","ð”¾"),y("gopf","ð•˜"),y("grave","`"),y("GreaterEqual","â‰¥"),y("GreaterEqualLess","â‹›"),y("GreaterFullEqual","â‰§"),y("GreaterGreater","âª¢"),y("GreaterLess","â‰·"),y("GreaterSlantEqual","â©¾"),y("GreaterTilde","â‰³"),y("Gscr","ð’¢"),y("gscr","â„Š"),y("gsim","â‰³"),y("gsime","âªŽ"),y("gsiml","âª"),y("gtcc","âª§"),y("gtcir","â©º"),y("gt",">"),y("GT",">"),y("Gt","â‰«"),y("gtdot","â‹—"),y("gtlPar","â¦•"),y("gtquest","â©¼"),y("gtrapprox","âª†"),y("gtrarr","â¥¸"),y("gtrdot","â‹—"),y("gtreqless","â‹›"),y("gtreqqless","âªŒ"),y("gtrless","â‰·"),y("gtrsim","â‰³"),y("gvertneqq","â‰©ï¸€"),y("gvnE","â‰©ï¸€"),y("Hacek","Ë‡"),y("hairsp","â€Š"),y("half","Â½"),y("hamilt","â„‹"),y("HARDcy","Ðª"),y("hardcy","ÑŠ"),y("harrcir","â¥ˆ"),y("harr","â†”"),y("hArr","â‡”"),y("harrw","â†­"),y("Hat","^"),y("hbar","â„"),y("Hcirc","Ä¤"),y("hcirc","Ä¥"),y("hearts","â™¥"),y("heartsuit","â™¥"),y("hellip","â€¦"),y("hercon","âŠ¹"),y("hfr","ð”¥"),y("Hfr","â„Œ"),y("HilbertSpace","â„‹"),y("hksearow","â¤¥"),y("hkswarow","â¤¦"),y("hoarr","â‡¿"),y("homtht","âˆ»"),y("hookleftarrow","â†©"),y("hookrightarrow","â†ª"),y("hopf","ð•™"),y("Hopf","â„"),y("horbar","â€•"),y("HorizontalLine","â”€"),y("hscr","ð’½"),y("Hscr","â„‹"),y("hslash","â„"),y("Hstrok","Ä¦"),y("hstrok","Ä§"),y("HumpDownHump","â‰Ž"),y("HumpEqual","â‰"),y("hybull","âƒ"),y("hyphen","â€"),y("Iacute","Ã"),y("iacute","Ã­"),y("ic","â£"),y("Icirc","ÃŽ"),y("icirc","Ã®"),y("Icy","Ð˜"),y("icy","Ð¸"),y("Idot","Ä°"),y("IEcy","Ð•"),y("iecy","Ðµ"),y("iexcl","Â¡"),y("iff","â‡”"),y("ifr","ð”¦"),y("Ifr","â„‘"),y("Igrave","ÃŒ"),y("igrave","Ã¬"),y("ii","â…ˆ"),y("iiiint","â¨Œ"),y("iiint","âˆ­"),y("iinfin","â§œ"),y("iiota","â„©"),y("IJlig","Ä²"),y("ijlig","Ä³"),y("Imacr","Äª"),y("imacr","Ä«"),y("image","â„‘"),y("ImaginaryI","â…ˆ"),y("imagline","â„"),y("imagpart","â„‘"),y("imath","Ä±"),y("Im","â„‘"),y("imof","âŠ·"),y("imped","Æµ"),y("Implies","â‡’"),y("incare","â„…"),y("in","âˆˆ"),y("infin","âˆž"),y("infintie","â§"),y("inodot","Ä±"),y("intcal","âŠº"),y("int","âˆ«"),y("Int","âˆ¬"),y("integers","â„¤"),y("Integral","âˆ«"),y("intercal","âŠº"),y("Intersection","â‹‚"),y("intlarhk","â¨—"),y("intprod","â¨¼"),y("InvisibleComma","â£"),y("InvisibleTimes","â¢"),y("IOcy","Ð"),y("iocy","Ñ‘"),y("Iogon","Ä®"),y("iogon","Ä¯"),y("Iopf","ð•€"),y("iopf","ð•š"),y("Iota","Î™"),y("iota","Î¹"),y("iprod","â¨¼"),y("iquest","Â¿"),y("iscr","ð’¾"),y("Iscr","â„"),y("isin","âˆˆ"),y("isindot","â‹µ"),y("isinE","â‹¹"),y("isins","â‹´"),y("isinsv","â‹³"),y("isinv","âˆˆ"),y("it","â¢"),y("Itilde","Ä¨"),y("itilde","Ä©"),y("Iukcy","Ð†"),y("iukcy","Ñ–"),y("Iuml","Ã"),y("iuml","Ã¯"),y("Jcirc","Ä´"),y("jcirc","Äµ"),y("Jcy","Ð™"),y("jcy","Ð¹"),y("Jfr","ð”"),y("jfr","ð”§"),y("jmath","È·"),y("Jopf","ð•"),y("jopf","ð•›"),y("Jscr","ð’¥"),y("jscr","ð’¿"),y("Jsercy","Ðˆ"),y("jsercy","Ñ˜"),y("Jukcy","Ð„"),y("jukcy","Ñ”"),y("Kappa","Îš"),y("kappa","Îº"),y("kappav","Ï°"),y("Kcedil","Ä¶"),y("kcedil","Ä·"),y("Kcy","Ðš"),y("kcy","Ðº"),y("Kfr","ð”Ž"),y("kfr","ð”¨"),y("kgreen","Ä¸"),y("KHcy","Ð¥"),y("khcy","Ñ…"),y("KJcy","ÐŒ"),y("kjcy","Ñœ"),y("Kopf","ð•‚"),y("kopf","ð•œ"),y("Kscr","ð’¦"),y("kscr","ð“€"),y("lAarr","â‡š"),y("Lacute","Ä¹"),y("lacute","Äº"),y("laemptyv","â¦´"),y("lagran","â„’"),y("Lambda","Î›"),y("lambda","Î»"),y("lang","âŸ¨"),y("Lang","âŸª"),y("langd","â¦‘"),y("langle","âŸ¨"),y("lap","âª…"),y("Laplacetrf","â„’"),y("laquo","Â«"),y("larrb","â‡¤"),y("larrbfs","â¤Ÿ"),y("larr","â†"),y("Larr","â†ž"),y("lArr","â‡"),y("larrfs","â¤"),y("larrhk","â†©"),y("larrlp","â†«"),y("larrpl","â¤¹"),y("larrsim","â¥³"),y("larrtl","â†¢"),y("latail","â¤™"),y("lAtail","â¤›"),y("lat","âª«"),y("late","âª­"),y("lates","âª­ï¸€"),y("lbarr","â¤Œ"),y("lBarr","â¤Ž"),y("lbbrk","â²"),y("lbrace","{"),y("lbrack","["),y("lbrke","â¦‹"),y("lbrksld","â¦"),y("lbrkslu","â¦"),y("Lcaron","Ä½"),y("lcaron","Ä¾"),y("Lcedil","Ä»"),y("lcedil","Ä¼"),y("lceil","âŒˆ"),y("lcub","{"),y("Lcy","Ð›"),y("lcy","Ð»"),y("ldca","â¤¶"),y("ldquo","â€œ"),y("ldquor","â€ž"),y("ldrdhar","â¥§"),y("ldrushar","â¥‹"),y("ldsh","â†²"),y("le","â‰¤"),y("lE","â‰¦"),y("LeftAngleBracket","âŸ¨"),y("LeftArrowBar","â‡¤"),y("leftarrow","â†"),y("LeftArrow","â†"),y("Leftarrow","â‡"),y("LeftArrowRightArrow","â‡†"),y("leftarrowtail","â†¢"),y("LeftCeiling","âŒˆ"),y("LeftDoubleBracket","âŸ¦"),y("LeftDownTeeVector","â¥¡"),y("LeftDownVectorBar","â¥™"),y("LeftDownVector","â‡ƒ"),y("LeftFloor","âŒŠ"),y("leftharpoondown","â†½"),y("leftharpoonup","â†¼"),y("leftleftarrows","â‡‡"),y("leftrightarrow","â†”"),y("LeftRightArrow","â†”"),y("Leftrightarrow","â‡”"),y("leftrightarrows","â‡†"),y("leftrightharpoons","â‡‹"),y("leftrightsquigarrow","â†­"),y("LeftRightVector","â¥Ž"),y("LeftTeeArrow","â†¤"),y("LeftTee","âŠ£"),y("LeftTeeVector","â¥š"),y("leftthreetimes","â‹‹"),y("LeftTriangleBar","â§"),y("LeftTriangle","âŠ²"),y("LeftTriangleEqual","âŠ´"),y("LeftUpDownVector","â¥‘"),y("LeftUpTeeVector","â¥ 
"),y("LeftUpVectorBar","â¥˜"),y("LeftUpVector","â†¿"),y("LeftVectorBar","â¥’"),y("LeftVector","â†¼"),y("lEg","âª‹"),y("leg","â‹š"),y("leq","â‰¤"),y("leqq","â‰¦"),y("leqslant","â©½"),y("lescc","âª¨"),y("les","â©½"),y("lesdot","â©¿"),y("lesdoto","âª"),y("lesdotor","âªƒ"),y("lesg","â‹šï¸€"),y("lesges","âª“"),y("lessapprox","âª…"),y("lessdot","â‹–"),y("lesseqgtr","â‹š"),y("lesseqqgtr","âª‹"),y("LessEqualGreater","â‹š"),y("LessFullEqual","â‰¦"),y("LessGreater","â‰¶"),y("lessgtr","â‰¶"),y("LessLess","âª¡"),y("lesssim","â‰²"),y("LessSlantEqual","â©½"),y("LessTilde","â‰²"),y("lfisht","â¥¼"),y("lfloor","âŒŠ"),y("Lfr","ð”"),y("lfr","ð”©"),y("lg","â‰¶"),y("lgE","âª‘"),y("lHar","â¥¢"),y("lhard","â†½"),y("lharu","â†¼"),y("lharul","â¥ª"),y("lhblk","â–„"),y("LJcy","Ð‰"),y("ljcy","Ñ™"),y("llarr","â‡‡"),y("ll","â‰ª"),y("Ll","â‹˜"),y("llcorner","âŒž"),y("Lleftarrow","â‡š"),y("llhard","â¥«"),y("lltri","â—º"),y("Lmidot","Ä¿"),y("lmidot","Å€"),y("lmoustache","âŽ°"),y("lmoust","âŽ°"),y("lnap","âª‰"),y("lnapprox","âª‰"),y("lne","âª‡"),y("lnE","â‰¨"),y("lneq","âª‡"),y("lneqq","â‰¨"),y("lnsim","â‹¦"),y("loang","âŸ¬"),y("loarr","â‡½"),y("lobrk","âŸ¦"),y("longleftarrow","âŸµ"),y("LongLeftArrow","âŸµ"),y("Longleftarrow","âŸ¸"),y("longleftrightarrow","âŸ·"),y("LongLeftRightArrow","âŸ·"),y("Longleftrightarrow","âŸº"),y("longmapsto","âŸ¼"),y("longrightarrow","âŸ¶"),y("LongRightArrow","âŸ¶"),y("Longrightarrow","âŸ¹"),y("looparrowleft","â†«"),y("looparrowright","â†¬"),y("lopar","â¦…"),y("Lopf","ð•ƒ"),y("lopf","ð•"),y("loplus","â¨­"),y("lotimes","â¨´"),y("lowast","âˆ—"),y("lowbar","_"),y("LowerLeftArrow","â†™"),y("LowerRightArrow","â†˜"),y("loz","â—Š"),y("lozenge","â—Š"),y("lozf","â§«"),y("lpar","("),y("lparlt","â¦“"),y("lrarr","â‡†"),y("lrcorner","âŒŸ"),y("lrhar","â‡‹"),y("lrhard","â¥­"),y("lrm","â€Ž"),y("lrtri","âŠ¿"),y("lsaquo","â€¹"),y("lscr","ð“"),y("Lscr","â„’"),y("lsh","â†°"),y("Lsh","â†°"),y("lsim","â‰²"),y("lsime","âª"),y("lsimg","âª"),y("lsqb","["),y("lsquo","â€˜"),y("lsquor","â€š"),y("Lstrok","Å"),y("lstrok","Å‚"),y("ltcc","âª¦"),y("ltcir","â©¹"),y("lt","<"),y("LT","<"),y("Lt","â‰ª"),y("ltdot","â‹–"),y("lthree","â‹‹"),y("ltimes","â‹‰"),y("ltlarr","â¥¶"),y("ltquest","â©»"),y("ltri","â—ƒ"),y("ltrie","âŠ´"),y("ltrif","â—‚"),y("ltrPar","â¦–"),y("lurdshar","â¥Š"),y("luruhar","â¥¦"),y("lvertneqq","â‰¨ï¸€"),y("lvnE","â‰¨ï¸€"),y("macr","Â¯"),y("male","â™‚"),y("malt","âœ 
"),y("maltese","âœ 
"),y("Map","â¤…"),y("map","â†¦"),y("mapsto","â†¦"),y("mapstodown","â†§"),y("mapstoleft","â†¤"),y("mapstoup","â†¥"),y("marker","â–®"),y("mcomma","â¨©"),y("Mcy","Ðœ"),y("mcy","Ð¼"),y("mdash","â€”"),y("mDDot","âˆº"),y("measuredangle","âˆ¡"),y("MediumSpace","âŸ"),y("Mellintrf","â„³"),y("Mfr","ð”"),y("mfr","ð”ª"),y("mho","â„§"),y("micro","Âµ"),y("midast","*"),y("midcir","â«°"),y("mid","âˆ£"),y("middot","Â·"),y("minusb","âŠŸ"),y("minus","âˆ’"),y("minusd","âˆ¸"),y("minusdu","â¨ª"),y("MinusPlus","âˆ“"),y("mlcp","â«›"),y("mldr","â€¦"),y("mnplus","âˆ“"),y("models","âŠ§"),y("Mopf","ð•„"),y("mopf","ð•ž"),y("mp","âˆ“"),y("mscr","ð“‚"),y("Mscr","â„³"),y("mstpos","âˆ¾"),y("Mu","Îœ"),y("mu","Î¼"),y("multimap","âŠ¸"),y("mumap","âŠ¸"),y("nabla","âˆ‡"),y("Nacute","Åƒ"),y("nacute","Å„"),y("nang","âˆ 
âƒ’"),y("nap","â‰‰"),y("napE","â©°Ì¸"),y("napid","â‰‹Ì¸"),y("napos","Å‰"),y("napprox","â‰‰"),y("natural","â™®"),y("naturals","â„•"),y("natur","â™®"),y("nbsp"," 
"),y("nbump","â‰ŽÌ¸"),y("nbumpe","â‰Ì¸"),y("ncap","â©ƒ"),y("Ncaron","Å‡"),y("ncaron","Åˆ"),y("Ncedil","Å…"),y("ncedil","Å†"),y("ncong","â‰‡"),y("ncongdot","â©­Ì¸"),y("ncup","â©‚"),y("Ncy","Ð"),y("ncy","Ð½"),y("ndash","â€“"),y("nearhk","â¤¤"),y("nearr","â†—"),y("neArr","â‡—"),y("nearrow","â†—"),y("ne","â‰ 
"),y("nedot","â‰Ì¸"),y("NegativeMediumSpace","â€‹"),y("NegativeThickSpace","â€‹"),y("NegativeThinSpace","â€‹"),y("NegativeVeryThinSpace","â€‹"),y("nequiv","â‰¢"),y("nesear","â¤¨"),y("nesim","â‰‚Ì¸"),y("NestedGreaterGreater","â‰«"),y("NestedLessLess","â‰ª"),y("NewLine","\n"),y("nexist","âˆ„"),y("nexists","âˆ„"),y("Nfr","ð”‘"),y("nfr","ð”«"),y("ngE","â‰§Ì¸"),y("nge","â‰±"),y("ngeq","â‰±"),y("ngeqq","â‰§Ì¸"),y("ngeqslant","â©¾Ì¸"),y("nges","â©¾Ì¸"),y("nGg","â‹™Ì¸"),y("ngsim","â‰µ"),y("nGt","â‰«âƒ’"),y("ngt","â‰¯"),y("ngtr","â‰¯"),y("nGtv","â‰«Ì¸"),y("nharr","â†®"),y("nhArr","â‡Ž"),y("nhpar","â«²"),y("ni","âˆ‹"),y("nis","â‹¼"),y("nisd","â‹º"),y("niv","âˆ‹"),y("NJcy","ÐŠ"),y("njcy","Ñš"),y("nlarr","â†š"),y("nlArr","â‡"),y("nldr","â€¥"),y("nlE","â‰¦Ì¸"),y("nle","â‰°"),y("nleftarrow","â†š"),y("nLeftarrow","â‡"),y("nleftrightarrow","â†®"),y("nLeftrightarrow","â‡Ž"),y("nleq","â‰°"),y("nleqq","â‰¦Ì¸"),y("nleqslant","â©½Ì¸"),y("nles","â©½Ì¸"),y("nless","â‰®"),y("nLl","â‹˜Ì¸"),y("nlsim","â‰´"),y("nLt","â‰ªâƒ’"),y("nlt","â‰®"),y("nltri","â‹ª"),y("nltrie","â‹¬"),y("nLtv","â‰ªÌ¸"),y("nmid","âˆ¤"),y("NoBreak","â 
"),y("NonBreakingSpace"," 
"),y("nopf","ð•Ÿ"),y("Nopf","â„•"),y("Not","â«¬"),y("not","Â¬"),y("NotCongruent","â‰¢"),y("NotCupCap","â‰­"),y("NotDoubleVerticalBar","âˆ¦"),y("NotElement","âˆ‰"),y("NotEqual","â‰ 
"),y("NotEqualTilde","â‰‚Ì¸"),y("NotExists","âˆ„"),y("NotGreater","â‰¯"),y("NotGreaterEqual","â‰±"),y("NotGreaterFullEqual","â‰§Ì¸"),y("NotGreaterGreater","â‰«Ì¸"),y("NotGreaterLess","â‰¹"),y("NotGreaterSlantEqual","â©¾Ì¸"),y("NotGreaterTilde","â‰µ"),y("NotHumpDownHump","â‰ŽÌ¸"),y("NotHumpEqual","â‰Ì¸"),y("notin","âˆ‰"),y("notindot","â‹µÌ¸"),y("notinE","â‹¹Ì¸"),y("notinva","âˆ‰"),y("notinvb","â‹·"),y("notinvc","â‹¶"),y("NotLeftTriangleBar","â§Ì¸"),y("NotLeftTriangle","â‹ª"),y("NotLeftTriangleEqual","â‹¬"),y("NotLess","â‰®"),y("NotLessEqual","â‰°"),y("NotLessGreater","â‰¸"),y("NotLessLess","â‰ªÌ¸"),y("NotLessSlantEqual","â©½Ì¸"),y("NotLessTilde","â‰´"),y("NotNestedGreaterGreater","âª¢Ì¸"),y("NotNestedLessLess","âª¡Ì¸"),y("notni","âˆŒ"),y("notniva","âˆŒ"),y("notnivb","â‹¾"),y("notnivc","â‹½"),y("NotPrecedes","âŠ€"),y("NotPrecedesEqual","âª¯Ì¸"),y("NotPrecedesSlantEqual","â‹ 
"),y("NotReverseElement","âˆŒ"),y("NotRightTriangleBar","â§Ì¸"),y("NotRightTriangle","â‹«"),y("NotRightTriangleEqual","â‹­"),y("NotSquareSubset","âŠÌ¸"),y("NotSquareSubsetEqual","â‹¢"),y("NotSquareSuperset","âŠÌ¸"),y("NotSquareSupersetEqual","â‹£"),y("NotSubset","âŠ‚âƒ’"),y("NotSubsetEqual","âŠˆ"),y("NotSucceeds","âŠ"),y("NotSucceedsEqual","âª°Ì¸"),y("NotSucceedsSlantEqual","â‹¡"),y("NotSucceedsTilde","â‰¿Ì¸"),y("NotSuperset","âŠƒâƒ’"),y("NotSupersetEqual","âŠ‰"),y("NotTilde","â‰"),y("NotTildeEqual","â‰„"),y("NotTildeFullEqual","â‰‡"),y("NotTildeTilde","â‰‰"),y("NotVerticalBar","âˆ¤"),y("nparallel","âˆ¦"),y("npar","âˆ¦"),y("nparsl","â«½âƒ¥"),y("npart","âˆ‚Ì¸"),y("npolint","â¨”"),y("npr","âŠ€"),y("nprcue","â‹ 
"),y("nprec","âŠ€"),y("npreceq","âª¯Ì¸"),y("npre","âª¯Ì¸"),y("nrarrc","â¤³Ì¸"),y("nrarr","â†›"),y("nrArr","â‡"),y("nrarrw","â†Ì¸"),y("nrightarrow","â†›"),y("nRightarrow","â‡"),y("nrtri","â‹«"),y("nrtrie","â‹­"),y("nsc","âŠ"),y("nsccue","â‹¡"),y("nsce","âª°Ì¸"),y("Nscr","ð’©"),y("nscr","ð“ƒ"),y("nshortmid","âˆ¤"),y("nshortparallel","âˆ¦"),y("nsim","â‰"),y("nsime","â‰„"),y("nsimeq","â‰„"),y("nsmid","âˆ¤"),y("nspar","âˆ¦"),y("nsqsube","â‹¢"),y("nsqsupe","â‹£"),y("nsub","âŠ„"),y("nsubE","â«…Ì¸"),y("nsube","âŠˆ"),y("nsubset","âŠ‚âƒ’"),y("nsubseteq","âŠˆ"),y("nsubseteqq","â«…Ì¸"),y("nsucc","âŠ"),y("nsucceq","âª°Ì¸"),y("nsup","âŠ…"),y("nsupE","â«†Ì¸"),y("nsupe","âŠ‰"),y("nsupset","âŠƒâƒ’"),y("nsupseteq","âŠ‰"),y("nsupseteqq","â«†Ì¸"),y("ntgl","â‰¹"),y("Ntilde","Ã‘"),y("ntilde","Ã±"),y("ntlg","â‰¸"),y("ntriangleleft","â‹ª"),y("ntrianglelefteq","â‹¬"),y("ntriangleright","â‹«"),y("ntrianglerighteq","â‹­"),y("Nu","Î"),y("nu","Î½"),y("num","#"),y("numero","â„–"),y("numsp","â€‡"),y("nvap","â‰âƒ’"),y("nvdash","âŠ¬"),y("nvDash","âŠ­"),y("nVdash","âŠ®"),y("nVDash","âŠ¯"),y("nvge","â‰¥âƒ’"),y("nvgt",">âƒ’"),y("nvHarr","â¤„"),y("nvinfin","â§ž"),y("nvlArr","â¤‚"),y("nvle","â‰¤âƒ’"),y("nvlt","<âƒ’"),y("nvltrie","âŠ´âƒ’"),y("nvrArr","â¤ƒ"),y("nvrtrie","âŠµâƒ’"),y("nvsim","âˆ¼âƒ’"),y("nwarhk","â¤£"),y("nwarr","â†–"),y("nwArr","â‡–"),y("nwarrow","â†–"),y("nwnear","â¤§"),y("Oacute","Ã“"),y("oacute","Ã³"),y("oast","âŠ›"),y("Ocirc","Ã”"),y("ocirc","Ã´"),y("ocir","âŠš"),y("Ocy","Ðž"),y("ocy","Ð¾"),y("odash","âŠ"),y("Odblac","Å"),y("odblac","Å‘"),y("odiv","â¨¸"),y("odot","âŠ™"),y("odsold","â¦¼"),y("OElig","Å’"),y("oelig","Å“"),y("ofcir","â¦¿"),y("Ofr","ð”’"),y("ofr","ð”¬"),y("ogon","Ë›"),y("Ograve","Ã’"),y("ograve","Ã²"),y("ogt","â§"),y("ohbar","â¦µ"),y("ohm","Î©"),y("oint","âˆ®"),y("olarr","â†º"),y("olcir","â¦¾"),y("olcross","â¦»"),y("oline","â€¾"),y("olt","â§€"),y("Omacr","ÅŒ"),y("omacr","Å"),y("Omega","Î©"),y("omega","Ï‰"),y("Omicron","ÎŸ"),y("omicron","Î¿"),y("omid","â¦¶"),y("ominus","âŠ–"),y("Oopf","ð•†"),y("oopf","ð• 
"),y("opar","â¦·"),y("OpenCurlyDoubleQuote","â€œ"),y("OpenCurlyQuote","â€˜"),y("operp","â¦¹"),y("oplus","âŠ•"),y("orarr","â†»"),y("Or","â©”"),y("or","âˆ¨"),y("ord","â©"),y("order","â„´"),y("orderof","â„´"),y("ordf","Âª"),y("ordm","Âº"),y("origof","âŠ¶"),y("oror","â©–"),y("orslope","â©—"),y("orv","â©›"),y("oS","â“ˆ"),y("Oscr","ð’ª"),y("oscr","â„´"),y("Oslash","Ã˜"),y("oslash","Ã¸"),y("osol","âŠ˜"),y("Otilde","Ã•"),y("otilde","Ãµ"),y("otimesas","â¨¶"),y("Otimes","â¨·"),y("otimes","âŠ—"),y("Ouml","Ã–"),y("ouml","Ã¶"),y("ovbar","âŒ½"),y("OverBar","â€¾"),y("OverBrace","âž"),y("OverBracket","âŽ´"),y("OverParenthesis","âœ"),y("para","Â¶"),y("parallel","âˆ¥"),y("par","âˆ¥"),y("parsim","â«³"),y("parsl","â«½"),y("part","âˆ‚"),y("PartialD","âˆ‚"),y("Pcy","ÐŸ"),y("pcy","Ð¿"),y("percnt","%"),y("period","."),y("permil","â€°"),y("perp","âŠ¥"),y("pertenk","â€±"),y("Pfr","ð”“"),y("pfr","ð”­"),y("Phi","Î¦"),y("phi","Ï†"),y("phiv","Ï•"),y("phmmat","â„³"),y("phone","â˜Ž"),y("Pi","Î 
"),y("pi","Ï€"),y("pitchfork","â‹”"),y("piv","Ï–"),y("planck","â„"),y("planckh","â„Ž"),y("plankv","â„"),y("plusacir","â¨£"),y("plusb","âŠž"),y("pluscir","â¨¢"),y("plus","+"),y("plusdo","âˆ”"),y("plusdu","â¨¥"),y("pluse","â©²"),y("PlusMinus","Â±"),y("plusmn","Â±"),y("plussim","â¨¦"),y("plustwo","â¨§"),y("pm","Â±"),y("Poincareplane","â„Œ"),y("pointint","â¨•"),y("popf","ð•¡"),y("Popf","â„™"),y("pound","Â£"),y("prap","âª·"),y("Pr","âª»"),y("pr","â‰º"),y("prcue","â‰¼"),y("precapprox","âª·"),y("prec","â‰º"),y("preccurlyeq","â‰¼"),y("Precedes","â‰º"),y("PrecedesEqual","âª¯"),y("PrecedesSlantEqual","â‰¼"),y("PrecedesTilde","â‰¾"),y("preceq","âª¯"),y("precnapprox","âª¹"),y("precneqq","âªµ"),y("precnsim","â‹¨"),y("pre","âª¯"),y("prE","âª³"),y("precsim","â‰¾"),y("prime","â€²"),y("Prime","â€³"),y("primes","â„™"),y("prnap","âª¹"),y("prnE","âªµ"),y("prnsim","â‹¨"),y("prod","âˆ"),y("Product","âˆ"),y("profalar","âŒ®"),y("profline","âŒ’"),y("profsurf","âŒ“"),y("prop","âˆ"),y("Proportional","âˆ"),y("Proportion","âˆ·"),y("propto","âˆ"),y("prsim","â‰¾"),y("prurel","âŠ°"),y("Pscr","ð’«"),y("pscr","ð“…"),y("Psi","Î¨"),y("psi","Ïˆ"),y("puncsp","â€ˆ"),y("Qfr","ð””"),y("qfr","ð”®"),y("qint","â¨Œ"),y("qopf","ð•¢"),y("Qopf","â„š"),y("qprime","â—"),y("Qscr","ð’¬"),y("qscr","ð“†"),y("quaternions","â„"),y("quatint","â¨–"),y("quest","?"),y("questeq","â‰Ÿ"),y("quot",'"'),y("QUOT",'"'),y("rAarr","â‡›"),y("race","âˆ½Ì±"),y("Racute","Å”"),y("racute","Å•"),y("radic","âˆš"),y("raemptyv","â¦³"),y("rang","âŸ©"),y("Rang","âŸ«"),y("rangd","â¦’"),y("range","â¦¥"),y("rangle","âŸ©"),y("raquo","Â»"),y("rarrap","â¥µ"),y("rarrb","â‡¥"),y("rarrbfs","â¤ 
"),y("rarrc","â¤³"),y("rarr","â†’"),y("Rarr","â† 
"),y("rArr","â‡’"),y("rarrfs","â¤ž"),y("rarrhk","â†ª"),y("rarrlp","â†¬"),y("rarrpl","â¥…"),y("rarrsim","â¥´"),y("Rarrtl","â¤–"),y("rarrtl","â†£"),y("rarrw","â†"),y("ratail","â¤š"),y("rAtail","â¤œ"),y("ratio","âˆ¶"),y("rationals","â„š"),y("rbarr","â¤"),y("rBarr","â¤"),y("RBarr","â¤"),y("rbbrk","â³"),y("rbrace","}"),y("rbrack","]"),y("rbrke","â¦Œ"),y("rbrksld","â¦Ž"),y("rbrkslu","â¦"),y("Rcaron","Å˜"),y("rcaron","Å™"),y("Rcedil","Å–"),y("rcedil","Å—"),y("rceil","âŒ‰"),y("rcub","}"),y("Rcy","Ð 
"),y("rcy","Ñ€"),y("rdca","â¤·"),y("rdldhar","â¥©"),y("rdquo","â€"),y("rdquor","â€"),y("rdsh","â†³"),y("real","â„œ"),y("realine","â„›"),y("realpart","â„œ"),y("reals","â„"),y("Re","â„œ"),y("rect","â–­"),y("reg","Â®"),y("REG","Â®"),y("ReverseElement","âˆ‹"),y("ReverseEquilibrium","â‡‹"),y("ReverseUpEquilibrium","â¥¯"),y("rfisht","â¥½"),y("rfloor","âŒ‹"),y("rfr","ð”¯"),y("Rfr","â„œ"),y("rHar","â¥¤"),y("rhard","â‡"),y("rharu","â‡€"),y("rharul","â¥¬"),y("Rho","Î¡"),y("rho","Ï"),y("rhov","Ï±"),y("RightAngleBracket","âŸ©"),y("RightArrowBar","â‡¥"),y("rightarrow","â†’"),y("RightArrow","â†’"),y("Rightarrow","â‡’"),y("RightArrowLeftArrow","â‡„"),y("rightarrowtail","â†£"),y("RightCeiling","âŒ‰"),y("RightDoubleBracket","âŸ§"),y("RightDownTeeVector","â¥"),y("RightDownVectorBar","â¥•"),y("RightDownVector","â‡‚"),y("RightFloor","âŒ‹"),y("rightharpoondown","â‡"),y("rightharpoonup","â‡€"),y("rightleftarrows","â‡„"),y("rightleftharpoons","â‡Œ"),y("rightrightarrows","â‡‰"),y("rightsquigarrow","â†"),y("RightTeeArrow","â†¦"),y("RightTee","âŠ¢"),y("RightTeeVector","â¥›"),y("rightthreetimes","â‹Œ"),y("RightTriangleBar","â§"),y("RightTriangle","âŠ³"),y("RightTriangleEqual","âŠµ"),y("RightUpDownVector","â¥"),y("RightUpTeeVector","â¥œ"),y("RightUpVectorBar","â¥”"),y("RightUpVector","â†¾"),y("RightVectorBar","â¥“"),y("RightVector","â‡€"),y("ring","Ëš"),y("risingdotseq","â‰“"),y("rlarr","â‡„"),y("rlhar","â‡Œ"),y("rlm","â€"),y("rmoustache","âŽ±"),y("rmoust","âŽ±"),y("rnmid","â«®"),y("roang","âŸ­"),y("roarr","â‡¾"),y("robrk","âŸ§"),y("ropar","â¦†"),y("ropf","ð•£"),y("Ropf","â„"),y("roplus","â¨®"),y("rotimes","â¨µ"),y("RoundImplies","â¥°"),y("rpar",")"),y("rpargt","â¦”"),y("rppolint","â¨’"),y("rrarr","â‡‰"),y("Rrightarrow","â‡›"),y("rsaquo","â€º"),y("rscr","ð“‡"),y("Rscr","â„›"),y("rsh","â†±"),y("Rsh","â†±"),y("rsqb","]"),y("rsquo","â€™"),y("rsquor","â€™"),y("rthree","â‹Œ"),y("rtimes","â‹Š"),y("rtri","â–¹"),y("rtrie","âŠµ"),y("rtrif","â–¸"),y("rtriltri","â§Ž"),y("RuleDelayed","â§´"),y("ruluhar","â¥¨"),y("rx","â„ž"),y("Sacute","Åš"),y("sacute","Å›"),y("sbquo","â€š"),y("scap","âª¸"),y("Scaron","Å 
"),y("scaron","Å¡"),y("Sc","âª¼"),y("sc","â‰»"),y("sccue","â‰½"),y("sce","âª°"),y("scE","âª´"),y("Scedil","Åž"),y("scedil","ÅŸ"),y("Scirc","Åœ"),y("scirc","Å"),y("scnap","âªº"),y("scnE","âª¶"),y("scnsim","â‹©"),y("scpolint","â¨“"),y("scsim","â‰¿"),y("Scy","Ð¡"),y("scy","Ñ"),y("sdotb","âŠ¡"),y("sdot","â‹…"),y("sdote","â©¦"),y("searhk","â¤¥"),y("searr","â†˜"),y("seArr","â‡˜"),y("searrow","â†˜"),y("sect","Â§"),y("semi",";"),y("seswar","â¤©"),y("setminus","âˆ–"),y("setmn","âˆ–"),y("sext","âœ¶"),y("Sfr","ð”–"),y("sfr","ð”°"),y("sfrown","âŒ¢"),y("sharp","â™¯"),y("SHCHcy","Ð©"),y("shchcy","Ñ‰"),y("SHcy","Ð¨"),y("shcy","Ñˆ"),y("ShortDownArrow","â†“"),y("ShortLeftArrow","â†"),y("shortmid","âˆ£"),y("shortparallel","âˆ¥"),y("ShortRightArrow","â†’"),y("ShortUpArrow","â†‘"),y("shy","Â­"),y("Sigma","Î£"),y("sigma","Ïƒ"),y("sigmaf","Ï‚"),y("sigmav","Ï‚"),y("sim","âˆ¼"),y("simdot","â©ª"),y("sime","â‰ƒ"),y("simeq","â‰ƒ"),y("simg","âªž"),y("simgE","âª 
"),y("siml","âª"),y("simlE","âªŸ"),y("simne","â‰†"),y("simplus","â¨¤"),y("simrarr","â¥²"),y("slarr","â†"),y("SmallCircle","âˆ˜"),y("smallsetminus","âˆ–"),y("smashp","â¨³"),y("smeparsl","â§¤"),y("smid","âˆ£"),y("smile","âŒ£"),y("smt","âªª"),y("smte","âª¬"),y("smtes","âª¬ï¸€"),y("SOFTcy","Ð¬"),y("softcy","ÑŒ"),y("solbar","âŒ¿"),y("solb","â§„"),y("sol","/"),y("Sopf","ð•Š"),y("sopf","ð•¤"),y("spades","â™ 
"),y("spadesuit","â™ 
"),y("spar","âˆ¥"),y("sqcap","âŠ“"),y("sqcaps","âŠ“ï¸€"),y("sqcup","âŠ”"),y("sqcups","âŠ”ï¸€"),y("Sqrt","âˆš"),y("sqsub","âŠ"),y("sqsube","âŠ‘"),y("sqsubset","âŠ"),y("sqsubseteq","âŠ‘"),y("sqsup","âŠ"),y("sqsupe","âŠ’"),y("sqsupset","âŠ"),y("sqsupseteq","âŠ’"),y("square","â–¡"),y("Square","â–¡"),y("SquareIntersection","âŠ“"),y("SquareSubset","âŠ"),y("SquareSubsetEqual","âŠ‘"),y("SquareSuperset","âŠ"),y("SquareSupersetEqual","âŠ’"),y("SquareUnion","âŠ”"),y("squarf","â–ª"),y("squ","â–¡"),y("squf","â–ª"),y("srarr","â†’"),y("Sscr","ð’®"),y("sscr","ð“ˆ"),y("ssetmn","âˆ–"),y("ssmile","âŒ£"),y("sstarf","â‹†"),y("Star","â‹†"),y("star","â˜†"),y("starf","â˜…"),y("straightepsilon","Ïµ"),y("straightphi","Ï•"),y("strns","Â¯"),y("sub","âŠ‚"),y("Sub","â‹"),y("subdot","âª½"),y("subE","â«…"),y("sube","âŠ†"),y("subedot","â«ƒ"),y("submult","â«"),y("subnE","â«‹"),y("subne","âŠŠ"),y("subplus","âª¿"),y("subrarr","â¥¹"),y("subset","âŠ‚"),y("Subset","â‹"),y("subseteq","âŠ†"),y("subseteqq","â«…"),y("SubsetEqual","âŠ†"),y("subsetneq","âŠŠ"),y("subsetneqq","â«‹"),y("subsim","â«‡"),y("subsub","â«•"),y("subsup","â«“"),y("succapprox","âª¸"),y("succ","â‰»"),y("succcurlyeq","â‰½"),y("Succeeds","â‰»"),y("SucceedsEqual","âª°"),y("SucceedsSlantEqual","â‰½"),y("SucceedsTilde","â‰¿"),y("succeq","âª°"),y("succnapprox","âªº"),y("succneqq","âª¶"),y("succnsim","â‹©"),y("succsim","â‰¿"),y("SuchThat","âˆ‹"),y("sum","âˆ‘"),y("Sum","âˆ‘"),y("sung","â™ª"),y("sup1","Â¹"),y("sup2","Â²"),y("sup3","Â³"),y("sup","âŠƒ"),y("Sup","â‹‘"),y("supdot","âª¾"),y("supdsub","â«˜"),y("supE","â«†"),y("supe","âŠ‡"),y("supedot","â«„"),y("Superset","âŠƒ"),y("SupersetEqual","âŠ‡"),y("suphsol","âŸ‰"),y("suphsub","â«—"),y("suplarr","â¥»"),y("supmult","â«‚"),y("supnE","â«Œ"),y("supne","âŠ‹"),y("supplus","â«€"),y("supset","âŠƒ"),y("Supset","â‹‘"),y("supseteq","âŠ‡"),y("supseteqq","â«†"),y("supsetneq","âŠ‹"),y("supsetneqq","â«Œ"),y("supsim","â«ˆ"),y("supsub","â«”"),y("supsup","â«–"),y("swarhk","â¤¦"),y("swarr","â†™"),y("swArr","â‡™"),y("swarrow","â†™"),y("swnwar","â¤ª"),y("szlig","ÃŸ"),y("Tab","\t"),y("target","âŒ–"),y("Tau","Î¤"),y("tau","Ï„"),y("tbrk","âŽ´"),y("Tcaron","Å¤"),y("tcaron","Å¥"),y("Tcedil","Å¢"),y("tcedil","Å£"),y("Tcy","Ð¢"),y("tcy","Ñ‚"),y("tdot","âƒ›"),y("telrec","âŒ•"),y("Tfr","ð”—"),y("tfr","ð”±"),y("there4","âˆ´"),y("therefore","âˆ´"),y("Therefore","âˆ´"),y("Theta","Î˜"),y("theta","Î¸"),y("thetasym","Ï‘"),y("thetav","Ï‘"),y("thickapprox","â‰ˆ"),y("thicksim","âˆ¼"),y("ThickSpace","âŸâ€Š"),y("ThinSpace","â€‰"),y("thinsp","â€‰"),y("thkap","â‰ˆ"),y("thksim","âˆ¼"),y("THORN","Ãž"),y("thorn","Ã¾"),y("tilde","Ëœ"),y("Tilde","âˆ¼"),y("TildeEqual","â‰ƒ"),y("TildeFullEqual","â‰…"),y("TildeTilde","â‰ˆ"),y("timesbar","â¨±"),y("timesb","âŠ 
"),y("times","Ã—"),y("timesd","â¨°"),y("tint","âˆ­"),y("toea","â¤¨"),y("topbot","âŒ¶"),y("topcir","â«±"),y("top","âŠ¤"),y("Topf","ð•‹"),y("topf","ð•¥"),y("topfork","â«š"),y("tosa","â¤©"),y("tprime","â€´"),y("trade","â„¢"),y("TRADE","â„¢"),y("triangle","â–µ"),y("triangledown","â–¿"),y("triangleleft","â—ƒ"),y("trianglelefteq","âŠ´"),y("triangleq","â‰œ"),y("triangleright","â–¹"),y("trianglerighteq","âŠµ"),y("tridot","â—¬"),y("trie","â‰œ"),y("triminus","â¨º"),y("TripleDot","âƒ›"),y("triplus","â¨¹"),y("trisb","â§"),y("tritime","â¨»"),y("trpezium","â¢"),y("Tscr","ð’¯"),y("tscr","ð“‰"),y("TScy","Ð¦"),y("tscy","Ñ†"),y("TSHcy","Ð‹"),y("tshcy","Ñ›"),y("Tstrok","Å¦"),y("tstrok","Å§"),y("twixt","â‰¬"),y("twoheadleftarrow","â†ž"),y("twoheadrightarrow","â† 
"),y("Uacute","Ãš"),y("uacute","Ãº"),y("uarr","â†‘"),y("Uarr","â†Ÿ"),y("uArr","â‡‘"),y("Uarrocir","â¥‰"),y("Ubrcy","ÐŽ"),y("ubrcy","Ñž"),y("Ubreve","Å¬"),y("ubreve","Å­"),y("Ucirc","Ã›"),y("ucirc","Ã»"),y("Ucy","Ð£"),y("ucy","Ñƒ"),y("udarr","â‡…"),y("Udblac","Å°"),y("udblac","Å±"),y("udhar","â¥®"),y("ufisht","â¥¾"),y("Ufr","ð”˜"),y("ufr","ð”²"),y("Ugrave","Ã™"),y("ugrave","Ã¹"),y("uHar","â¥£"),y("uharl","â†¿"),y("uharr","â†¾"),y("uhblk","â–€"),y("ulcorn","âŒœ"),y("ulcorner","âŒœ"),y("ulcrop","âŒ"),y("ultri","â—¸"),y("Umacr","Åª"),y("umacr","Å«"),y("uml","Â¨"),y("UnderBar","_"),y("UnderBrace","âŸ"),y("UnderBracket","âŽµ"),y("UnderParenthesis","â"),y("Union","â‹ƒ"),y("UnionPlus","âŠŽ"),y("Uogon","Å²"),y("uogon","Å³"),y("Uopf","ð•Œ"),y("uopf","ð•¦"),y("UpArrowBar","â¤’"),y("uparrow","â†‘"),y("UpArrow","â†‘"),y("Uparrow","â‡‘"),y("UpArrowDownArrow","â‡…"),y("updownarrow","â†•"),y("UpDownArrow","â†•"),y("Updownarrow","â‡•"),y("UpEquilibrium","â¥®"),y("upharpoonleft","â†¿"),y("upharpoonright","â†¾"),y("uplus","âŠŽ"),y("UpperLeftArrow","â†–"),y("UpperRightArrow","â†—"),y("upsi","Ï…"),y("Upsi","Ï’"),y("upsih","Ï’"),y("Upsilon","Î¥"),y("upsilon","Ï…"),y("UpTeeArrow","â†¥"),y("UpTee","âŠ¥"),y("upuparrows","â‡ˆ"),y("urcorn","âŒ"),y("urcorner","âŒ"),y("urcrop","âŒŽ"),y("Uring","Å®"),y("uring","Å¯"),y("urtri","â—¹"),y("Uscr","ð’°"),y("uscr","ð“Š"),y("utdot","â‹°"),y("Utilde","Å¨"),y("utilde","Å©"),y("utri","â–µ"),y("utrif","â–´"),y("uuarr","â‡ˆ"),y("Uuml","Ãœ"),y("uuml","Ã¼"),y("uwangle","â¦§"),y("vangrt","â¦œ"),y("varepsilon","Ïµ"),y("varkappa","Ï°"),y("varnothing","âˆ…"),y("varphi","Ï•"),y("varpi","Ï–"),y("varpropto","âˆ"),y("varr","â†•"),y("vArr","â‡•"),y("varrho","Ï±"),y("varsigma","Ï‚"),y("varsubsetneq","âŠŠï¸€"),y("varsubsetneqq","â«‹ï¸€"),y("varsupsetneq","âŠ‹ï¸€"),y("varsupsetneqq","â«Œï¸€"),y("vartheta","Ï‘"),y("vartriangleleft","âŠ²"),y("vartriangleright","âŠ³"),y("vBar","â«¨"),y("Vbar","â««"),y("vBarv","â«©"),y("Vcy","Ð’"),y("vcy","Ð²"),y("vdash","âŠ¢"),y("vDash","âŠ¨"),y("Vdash","âŠ©"),y("VDash","âŠ«"),y("Vdashl","â«¦"),y("veebar","âŠ»"),y("vee","âˆ¨"),y("Vee","â‹"),y("veeeq","â‰š"),y("vellip","â‹®"),y("verbar","|"),y("Verbar","â€–"),y("vert","|"),y("Vert","â€–"),y("VerticalBar","âˆ£"),y("VerticalLine","|"),y("VerticalSeparator","â˜"),y("VerticalTilde","â‰€"),y("VeryThinSpace","â€Š"),y("Vfr","ð”™"),y("vfr","ð”³"),y("vltri","âŠ²"),y("vnsub","âŠ‚âƒ’"),y("vnsup","âŠƒâƒ’"),y("Vopf","ð•"),y("vopf","ð•§"),y("vprop","âˆ"),y("vrtri","âŠ³"),y("Vscr","ð’±"),y("vscr","ð“‹"),y("vsubnE","â«‹ï¸€"),y("vsubne","âŠŠï¸€"),y("vsupnE","â«Œï¸€"),y("vsupne","âŠ‹ï¸€"),y("Vvdash","âŠª"),y("vzigzag","â¦š"),y("Wcirc","Å´"),y("wcirc","Åµ"),y("wedbar","â©Ÿ"),y("wedge","âˆ§"),y("Wedge","â‹€"),y("wedgeq","â‰™"),y("weierp","â„˜"),y("Wfr","ð”š"),y("wfr","ð”´"),y("Wopf","ð•Ž"),y("wopf","ð•¨"),y("wp","â„˜"),y("wr","â‰€"),y("wreath","â‰€"),y("Wscr","ð’²"),y("wscr","ð“Œ"),y("xcap","â‹‚"),y("xcirc","â—¯"),y("xcup","â‹ƒ"),y("xdtri","â–½"),y("Xfr","ð”›"),y("xfr","ð”µ"),y("xharr","âŸ·"),y("xhArr","âŸº"),y("Xi","Îž"),y("xi","Î¾"),y("xlarr","âŸµ"),y("xlArr","âŸ¸"),y("xmap","âŸ¼"),y("xnis","â‹»"),y("xodot","â¨€"),y("Xopf","ð•"),y("xopf","ð•©"),y("xoplus","â¨"),y("xotime","â¨‚"),y("xrarr","âŸ¶"),y("xrArr","âŸ¹"),y("Xscr","ð’³"),y("xscr","ð“"),y("xsqcup","â¨†"),y("xuplus","â¨„"),y("xutri","â–³"),y("xvee","â‹"),y("xwedge","â‹€"),y("Yacute","Ã"),y("yacute","Ã½"),y("YAcy","Ð¯"),y("yacy","Ñ"),y("Ycirc","Å¶"),y("ycirc","Å·"),y("Ycy","Ð«"),y("ycy","Ñ‹"),y("yen","Â¥"),y("Yfr","ð”œ"),y("yfr","ð”¶"),y("YIcy","Ð‡"),y("yicy","Ñ—"),y("Yopf","ð•"),y("yopf","ð•ª"),y("Yscr","ð’´"),y("yscr","ð“Ž"),y("YUcy","Ð®"),y("yucy","ÑŽ"),y("yuml","Ã¿"),y("Yuml","Å¸"),y("Zacute","Å¹"),y("zacute","Åº"),y("Zcaron","Å½"),y("zcaron","Å¾"),y("Zcy","Ð—"),y("zcy","Ð·"),y("Zdot","Å»"),y("zdot","Å¼"),y("zeetrf","â„¨"),y("ZeroWidthSpace","â€‹"),y("Zeta","Î–"),y("zeta","Î¶"),y("zfr","ð”·"),y("Zfr","â„¨"),y("ZHcy","Ð–"),y("zhcy","Ð¶"),y("zigrarr","â‡"),y("zopf","ð•«"),y("Zopf","â„¤"),y("Zscr","ð’µ"),y("zscr","ð“"),y("zwj","â€"),y("zwnj","â€Œ")])),Dc=f(ac,function(r){return 
f(ju,"&"+r+";",f(xu,r,Lc))},Ua(Ca(Bt))),Tc=e(function(r,n,t){return 
n(r(t))}),Nc=M,Sc=function(r){return f(Nc,r,"")},Cc=function(r){return 
x(0>r||r>1114111?"ï¿½":r>65535?String.fromCharCode(Math.floor((r-=65536)/1024)+55296,r%1024+56320):String.fromCharCode(r))},Vc=j,Rc=e(function(r,n,t){r:for(;;){if(!n.b)return 
pt(t);var e=n.a,u=n.b;switch(e){case"0":r=a=r-1,n=c=u,t=i=t;continue 
r;case"1":var a=r-1,c=u,i=t+f(Vc,16,r);r=a,n=c,t=i;continue 
r;case"2":a=r-1,c=u,i=t+2*f(Vc,16,r),r=a,n=c,t=i;continue 
r;case"3":a=r-1,c=u,i=t+3*f(Vc,16,r),r=a,n=c,t=i;continue 
r;case"4":a=r-1,c=u,i=t+4*f(Vc,16,r),r=a,n=c,t=i;continue 
r;case"5":a=r-1,c=u,i=t+5*f(Vc,16,r),r=a,n=c,t=i;continue 
r;case"6":a=r-1,c=u,i=t+6*f(Vc,16,r),r=a,n=c,t=i;continue 
r;case"7":a=r-1,c=u,i=t+7*f(Vc,16,r),r=a,n=c,t=i;continue 
r;case"8":a=r-1,c=u,i=t+8*f(Vc,16,r),r=a,n=c,t=i;continue 
r;case"9":a=r-1,c=u,i=t+9*f(Vc,16,r),r=a,n=c,t=i;continue 
r;case"a":a=r-1,c=u,i=t+10*f(Vc,16,r),r=a,n=c,t=i;continue 
r;case"b":a=r-1,c=u,i=t+11*f(Vc,16,r),r=a,n=c,t=i;continue 
r;case"c":a=r-1,c=u,i=t+12*f(Vc,16,r),r=a,n=c,t=i;continue 
r;case"d":a=r-1,c=u,i=t+13*f(Vc,16,r),r=a,n=c,t=i;continue 
r;case"e":a=r-1,c=u,i=t+14*f(Vc,16,r),r=a,n=c,t=i;continue 
r;case"f":a=r-1,c=u,i=t+15*f(Vc,16,r),r=a,n=c,t=i;continue 
r;default:return lt(Sc(e)+" is not a valid hexadecimal 
character.")}}}),Uc=t(function(r,n){return 
n.$?lt(n.a):pt(r(n.a))}),Bc=t(function(r,n){return 
n.$?lt(r(n.a)):pt(n.a)}),Oc=K,Pc=function(r){return 
s(Oc,ot,L,r)},jc=function(r){if(we(r))return lt("Empty strings are not 
valid hexadecimal strings.");var n=function(){if(f(xe,"-",r)){var 
n=f(ju,L,function(r){return r.b?ht(r.b):$t}(Pc(r)));return 
f(Uc,qa,s(Rc,Dt(n)-1,n,0))}return s(Rc,ge(r)-1,Pc(r),0)}();return 
f(Bc,function(n){return f(At," ",N(['"'+r+'"',"is not a valid hexadecimal 
string because",n]))},n)},Ic=function(r){var n=Vt(r);return 
n>=48&&57>=n||n>=65&&70>=n||n>=97&&102>=n},_c=f(da,function(r){var 
n=jc(za(r));return 
n.$?Ia(n.a):Ga(n.a)},Ua(Ca(Ic))),Gc={$:1},zc=Yn,Hc=Qn,Mc=t(function(r,n){return{bB:n.bB+(r-n.b),d:n.d,e:n.e,b:r,dP:n.dP,cy:n.cy}}),Fc=Xn,Kc=Jn,Zc=t(function(r,n){if(s(Kc,101,r,n)||s(Kc,69,r,n)){var 
t=r+1,e=s(Kc,43,t,n)||s(Kc,45,t,n)?t+1:t,u=f(Fc,e,n);return 
p(e,u)?-u:u}return r}),Jc=t(function(r,n){return 
f(Zc,s(Kc,46,r,n)?f(Fc,r+1,n):r,n)}),Xc=a(function(r,n,t,e,u){var 
a=e.a,c=e.b;if(1===n.$)return f(sa,!0,f(wa,u,n.a));var i=n.a;return 
p(t,a)?f(sa,0>h(u.b,t),f(wa,u,r)):s(la,!0,i(c),f(Mc,a,u))}),Yc=function(r){if(0===r.length||/[\sxbo]/.test(r))return 
$t;var n=+r;return n==n?ht(n):$t},Qc=c(function(r,n,t,e,u,a){var 
c=u.a,i=f(Jc,c,a.cy);if(0>i)return 
f(sa,!0,l(Ka,a.dP,a.bB-(i+a.b),r,a.d));if(p(a.b,i))return 
f(sa,!1,f(wa,a,n));if(p(c,i))return b(Xc,r,t,a.b,u,a);if(1===e.$)return 
f(sa,!0,f(wa,a,r));var o=e.a,d=Yc(s(he,a.b,i,a.cy));return 
1===d.$?f(sa,!0,f(wa,a,r)):s(la,!0,o(d.a),f(Mc,i,a))}),Wc=f(t(function(r,n){return 
t={bx:lt(n),bM:r,bN:lt(n),bS:lt(n),bX:pt(pe),dh:n,b8:lt(n)},function(r){if(s(Kc,48,r.b,r.cy)){var 
n=r.b+1,e=n+1;return 
s(Kc,120,n,r.cy)?b(Xc,t.dh,t.bS,e,f(Hc,e,r.cy),r):s(Kc,111,n,r.cy)?b(Xc,t.dh,t.b8,e,s(zc,8,e,r.cy),r):s(Kc,98,n,r.cy)?b(Xc,t.dh,t.bx,e,s(zc,2,e,r.cy),r):d(Qc,t.dh,t.bM,t.bX,t.bN,y(n,0),r)}return 
d(Qc,t.dh,t.bM,t.bX,t.bN,s(zc,10,r.b,r.cy),r)};var 
t}),Gc,Gc),ri=(Be=qc(N([f(Wa,f(Sa,Ga(pe),Aa(function(r){return"x"===r||"X"===r})),_c),f(Wa,f(Sa,Ga(pe),La(Qt("0"))),Wc)])),f(Wa,f(Sa,Ga(pe),Aa(Qt("#"))),f(ac,f(Tc,Cc,Sc),Be))),ni=f(Wa,f(Sa,Ga(pe),Aa(Qt("&"))),qc(N([f(Sa,pa(Dc),kc),f(Sa,pa(ri),kc),Ga("&")]))),ti=function(r){return 
f(Wa,f(Sa,Ga(pe),Aa(Qt(r))),f(Sa,f(ac,At(""),xc(qc(N([Ua(Ca(function(n){return!p(n,r)&&"&"!==n})),ni])))),Aa(Qt(r))))},ei=function(r){return!r.b},ui=t(function(r,n){return 
f(hc,L,function(t){return qc(N([f(ac,function(r){return 
bc(f(ot,r,t))},n),ei(t)?Ia("expecting at least one 
"+r):Ga(lc(jt(t)))]))})}),ai=f(ac,At(""),f(ui,"attribute 
value",qc(N([Ua(Ca(function(r){return!Ba(r)&&'"'!==r&&"'"!==r&&"="!==r&&"<"!==r&&">"!==r&&"`"!==r&&"&"!==r})),ni])))),ci=qc(N([f(Wa,f(Sa,f(Sa,Ga(pe),Aa(Qt("="))),La(Ba)),qc(N([ai,ti('"'),ti("'")]))),Ga("")])),ii=f(Wa,f(Wa,Ga(We),f(Sa,Ac,La(Ba))),f(Sa,ci,La(Ba))),oi=xc(ii),fi=f(ac,za,Ua(f(Sa,Aa(Pt),La(function(r){return 
Pt(r)||"-"===r})))),si=function(r){return{$:0,a:r}},li=f(ac,f(Tc,At(""),si),f(ui,"text 
element",qc(N([Ua(Ca(function(r){return"<"!==r&&"&"!==r})),ni]))));function 
bi(){return qc(N([li,cc,di()]))}function di(){return f(da,function(r){var 
n=r.a,t=r.b;return 
sc(n)?f(Sa,f(Sa,Ga(s(fa,n,t,L)),qc(N([Aa(Qt("/")),Ga(0)]))),Aa(Qt(">"))):f(Wa,f(Sa,Ga(f(fa,n,t)),Aa(Qt(">"))),f(Sa,xc(pa(bi())),Ha(n)))},f(Wa,f(Wa,f(Sa,Ga(We),Aa(Qt("<"))),f(Sa,fi,La(Ba))),oi))}var 
vi=bi();bi=function(){return vi};var pi=di();di=function(){return pi};var 
mi,gi,hi=e(function(r,n,t){return{bB:n,dH:t,dP:r}}),$i=function(r){return 
s(hi,r.dP,r.bB,r.dH)},wi=t(function(r,n){r:for(;;)switch(r.$){case 
0:return n;case 1:var t=r.b;r=r.a,n=f(ot,t,n);continue r;default:var 
e=r.b;r=r.a,n=f(wi,e,n);continue r}}),yi=t(function(r,n){var 
t=r({bB:1,d:L,e:1,b:0,dP:1,cy:n});return 
t.$?lt(f(wi,t.b,L)):pt(t.b)}),qi=t(function(r,n){var t=f(yi,r,n);return 
t.$?lt(f(Te,$i,t.a)):pt(t.a)}),xi=f(Je,function(r){return 
1===r.$?oa:de({$:1,a:r.a})},f(le,function(r){return 
we(r)?pt(L):f(qi,f(ui,"node",vi),r)},f(Xe,"formatted_body",ru))),Ai=f(Je,function(r){return 
r.$||"org.matrix.custom.html"!==r.a?oa:xi},Yu(f(Xe,"format",ru))),ki=e(function(r,n,t){return{cV:r,a1:t,ef:n}}),Ei=s(be,t(function(r,n){return{aB:n,aQ:r}}),f(Xe,"w",Ye),f(Xe,"h",Ye)),Li=dr,Di=b(Li,u(function(r,n,t,e){return{aB:n,cE:e,cF:t,aQ:r}}),f(Xe,"w",Ye),f(Xe,"h",Ye),Yu(f(Xe,"thumbnail_url",ru)),Yu(f(Xe,"thumbnail_info",Ei))),Ti=l(Xu,ki,f(Xe,"body",ru),f(Xe,"url",ru),Yu(f(Xe,"info",Di))),Ni=s(be,t(function(r,n){return{cV:r,ef:n}}),f(Xe,"body",ru),f(Xe,"url",ru)),Si=f(Je,function(r){switch(r){case"m.text":return 
f(le,ua,Ai);case"m.emote":return f(le,ra,Ai);case"m.notice":return 
f(le,ea,Ai);case"m.image":return f(le,ta,Ti);case"m.file":return 
f(le,na,ia);case"m.audio":return f(le,Wu,ca);case"m.video":return 
f(le,aa,Ni);default:return Zu("Unsupported msgtype: 
"+r)}},f(Xe,"msgtype",ru)),Ci=pe,Vi=t(function(r,n){return 
b(Li,r,f(Xe,"type",ru),f(Xe,"content",n),f(Xe,"sender",ru),f(le,Ci,f(Xe,"origin_server_ts",Ye)))}),Ri=Qe(N([f(Je,function(r){switch(r){case"m.room.message":return 
f(Vi,u(function(r,n,t,e){return{$:0,a:{ae:n,aA:r,_:e,aL:t}}}),Si);case"m.room.member":return 
f(Je,function(r){return f(Vi,u(function(n,t,e,u){return 
f(Ku,r,{ae:t,aA:n,_:u,aL:e})}),Qu)},f(Xe,"state_key",ru));default:return 
Zu("Unsupported event type: 
"+r)}},f(Xe,"type",ru)),f(Vi,u(function(r,n,t,e){return{$:2,a:{ae:n,aA:r,_:e,aL:t}}}),de(0))])),Ui=function(r){return{$:3,b:r}},Bi=l(Xu,e(function(r,n,t){return{bA:t,bK:n,bm:r}}),f(Xe,"start",ru),Yu(f(Xe,"end",ru)),f(Xe,"chunk",Ui(Ri))),Oi={$:0},Pi=t(function(r,n){return 
f(Fu,r,{cV:Oi,dm:"GET",dF:L,dG:N(["rooms",n,"initialSync"]),dM:f(Xe,"messages",Bi)})}),ji=s(be,t(function(r,n){return 
y(r,n)}),f(Xe,"state_key",ru),f(Xe,"content",Qu)),Ii=f(Xe,"chunk",f(Je,f(Tc,Ec,de),Ui(ji))),_i=t(function(r,n){return 
f(Fu,r,{cV:Oi,dm:"GET",dF:L,dG:N(["rooms",n,"members"]),dM:Ii})}),Gi=t(function(r,n){return 
f(Fu,r,{cV:Oi,dm:"GET",dF:L,dG:N(["directory","room",n]),dM:f(le,pe,f(Xe,"room_id",ru))})}),zi=function(r){return 
r},Hi=V,Mi=function(r){return f(Hi,function(r){switch(r.$){case 0:return 
zi(r.a._);case 1:return zi(r.b._);default:return 
zi(r.a._)}},r)},Fi=t(function(r,n){return f(Se,pe,f(Ne,function(n){return 
f(Se,function(r){return{bK:n.bK,C:n.C,b2:r,dO:n.dO,F:n.F,bm:n.bm}},f(_i,r,n.F))},f(Ne,function(n){return 
f(Se,function(r){return{bK:(t=r.bK,1===t.$?r.bm:t.a),C:Mi(r.bA),dO:n.dO,F:n.F,bm:r.bm};var 
t},f(Pi,r,n.F))},f(Se,function(r){return{dO:n,F:r}},f(Gi,r,n)))))}),Ki={V:"",ai:""},Zi=t(function(r,n){return{$:0,a:r,b:n}}),Ji=$u,Xi={$:10},Yi=function(r){return 
function(n){return 
p(ge(n.cy),n.b)?s(la,!1,0,n):f(sa,!1,f(wa,n,r))}},Qi=Yi(Xi),Wi=function(r){return 
Rt(r)||Ot(r)||f(oc,r,N(["-",".","=","_","/"]))},ro=function(r){return 
Pt(r)||f(oc,r,N([".","-",":"]))},no={$:7},to=t(function(r,n){return!f(xu,r,n).$}),eo=t(function(r,n){return 
f(to,r,n)}),uo=i(function(r,n,t,e,u,a,c){for(;;){var 
i=s(ya,r,n,u);if(p(i,-1))return{bB:e,d:c,e:a,b:n,dP:t,cy:u};p(i,-2)?(r=r,n+=1,t+=1,e=1,u=u,a=a,c=c):(r=r,n=i,t=t,e+=1,u=u,a=a,c=c)}}),ao=function(r){return 
function(r){return function(n){var t=s(ya,r.bm,n.b,n.cy);if(p(t,-1))return 
f(sa,!1,f(wa,n,r.bM));var 
e=p(t,-2)?v(uo,r.a2,n.b+1,n.dP+1,1,n.cy,n.e,n.d):v(uo,r.a2,t,n.dP,n.bB+1,n.cy,n.e,n.d),u=s(he,n.b,e.b,n.cy);return 
f(eo,u,r.bc)?f(sa,!1,f(wa,n,r.bM)):s(la,!0,u,e)}}({bM:no,a2:r.a2,bc:r.bc,bm:r.bm})},co=f(Wa,f(Wa,f(Sa,Ga(Zi),tc("@")),f(Sa,ao({a2:Wi,bc:Ji,bm:Wi}),tc(":"))),f(Sa,ao({a2:ro,bc:Ji,bm:ro}),Qi)),io=f(Tc,za,f(Tc,qi(co),Bc(function(){return"Must 
follow format: 
@user:example.com"}))),oo={Q:$t,a5:$t,am:"",I:2,aO:f(Bc,function(){return"Something's 
wrong with the user ID 
parser"},io("@alice:example.com")),bs:""},fo=Fn,so=t(function(r,n){return 
f(Fu,r,{cV:f(fo,"application/json","{}"),dm:"POST",dF:L,dG:N(["join",n]),dM:de(0)})}),lo=function(r){return 
r.aC},bo=t(function(r,n){return 
1===lo(r)?f(so,r,n):ke(0)}),vo=Fe(L),po=t(function(r,n){return{$:0,a:r,b:n}}),mo=(mi=Ci,Tr(function(r){r(Dr(mi(Date.now())))})),go=o(function(r,n,t,e,u,a,c,i,o){return{bC:e,aR:r,aV:i,a4:c,ba:a,cs:n,cu:t,cz:u,br:o}}),ho=Qe(N([f(le,ht,er),f(le,ht,f(Je,function(r){switch(za(r)){case"true":return 
de(!0);case"false":return de(!1);default:return 
Zu("")}},ru))])),$o=f(Je,function(r){return f(me," 
",r)?Zu("commentSectionId can't contain 
spaces"):de(r)},f(Je,function(r){return f(me,"_",r)?Zu("commentSectionId 
can't contain 
underscores"):de(r)},ru)),wo=Qe(N([f(le,ht,ur),f(le,Yc,ru)])),yo=Qe(N([f(le,ht,Ye),f(le,qe,ru)])),qo=a(function(r,n,t,e,u){return{U:u,bT:r,aC:n,aq:t,aO:e}}),xo=f(Je,function(r){switch(r){case"guest":return 
de(0);case"user":return de(1);default:return Zu(r+" is not a valid 
Session.Kind")}},ru),Ao=f(le,pe,d(vr,qo,f(Xe,"homeserverUrl",ru),f(Xe,"kind",xo),f(Xe,"txnId",Ye),f(Xe,"userId",ru),f(Xe,"accessToken",ru))),ko=function(r){return{$:5,c:r}},Eo=be(Yt),Lo=mr,Do=ar,To=e(function(r,n,t){return 
f(Je,function(e){var u=f(Lo,r,e);if(u.$)return de(t);var 
a=u.a,c=f(Lo,Qe(N([n,ko(t)])),a);return 
c.$?Zu(_t(c.a)):de(c.a)},Do)}),No=u(function(r,n,t,e){return 
f(Eo,s(To,f(Xe,r,Do),n,t),e)}),So=e(function(r,n,t){return 
f(Eo,f(Xe,r,n),t)}),Co=l(No,"updateInterval",wo,$t,l(No,"guestPostingEnabled",ho,$t,l(No,"loginEnabled",ho,$t,l(No,"pageSize",yo,$t,l(No,"storedSession",(gi=Ao,Qe(N([ko($t),f(le,ht,gi)]))),$t,s(So,"commentSectionId",$o,s(So,"siteName",ru,s(So,"serverName",ru,s(So,"defaultHomeserverUrl",ru,de(go)))))))))),Vo=c(function(r,n,t,e,u,a){return{aR:r,aV:u,a4:e,ba:t,dO:n,br:a}}),Ro=function(r){return 
r},Uo=function(r){return 
y(d(Vo,r.aR,"#comments_"+(n=r).cu+"_"+n.bC+":"+n.cs,f(ju,10,r.ba),f(ju,!0,r.a4),f(ju,!0,r.aV),Ro(f(ju,0,r.br))),r.cz);var 
n},Bo=t(function(r,n){return 
f(le,pe,s(be,s(qo,r,n,0),f(Xe,"user_id",ru),f(Xe,"access_token",ru)))}),Oo=t(function(r,n){return{$:0,a:r,b:n}}),Po=t(function(r,n){return 
f(Oo,zu(r),zu(n))}),jo=function(r){return 
Iu({U:$t,cV:r.cV,dm:r.dm,dM:r.dM,ef:r.ef})},Io=function(r){return 
jo({cV:f(fo,"application/json","{}"),dm:"POST",dM:f(Bo,r,0),ef:s(Mu,r,N(["register"]),N([f(Po,"kind","guest")]))})},_o=t(function(r,n){return{$:0,a:r,b:n}}),Go=t(function(r,n){return{cg:n,cD:r}}),zo=ke(f(Go,$u,$u)),Ho=t(function(r,n){var 
t=r.a,e=r.b,u=f(xu,t,n);return 
s(Lu,t,1===u.$?N([e]):f(ot,e,u.a),n)}),Mo=e(function(r,n,t){for(;;){if(-2===t.$)return 
n;var 
e=t.e,u=r,a=s(r,t.b,t.c,s(Mo,r,n,t.d));r=u,n=a,t=e}}),Fo=c(function(r,n,u,a,c,i){var 
o=s(Mo,e(function(t,e,a){r:for(;;){var c=a.a,i=a.b;if(c.b){var 
o=c.a,f=o.a,b=o.b,d=c.b;if(0>h(f,t)){t=t,e=e,a=y(d,s(r,f,b,i));continue 
r}return h(f,t)>0?y(c,s(u,t,e,i)):y(d,l(n,f,b,e,i))}return 
y(c,s(u,t,e,i))}}),y(st(a),i),c),f=o.a,b=o.b;return 
s(Lt,t(function(n,t){return 
s(r,n.a,n.b,t)}),b,f)}),Ko=rt,Zo=Rr,Jo=e(function(r,n,t){if(n.b){var 
e=n.a,u=n.b,a=Zo(f(Ko,e,f(yu,r,e)));return f(Ne,function(n){return 
s(Jo,r,u,s(Lu,e,n,t))},a)}return 
ke(t)});Gr.Time=zr(zo,e(function(r,n,t){var a=t.cg,c=e(function(r,n,t){var 
e,u=t.c;return q(t.a,t.b,f(Ne,function(){return u},(e=n,Tr(function(r){var 
n=e.f;2===n.$&&n.c&&n.c(),e.f=null,r(Dr(w))}))))}),i=s(Lt,Ho,$u,n),o=d(Fo,e(function(r,n,t){var 
e=t.b,u=t.c;return q(f(ot,r,t.a),e,u)}),u(function(r,n,t,e){var 
u=e.c;return 
q(e.a,s(Lu,r,t,e.b),u)}),c,i,a,q(L,$u,ke(0))),l=o.a,b=o.b;return 
f(Ne,function(r){return ke(f(Go,i,r))},f(Ne,function(){return 
s(Jo,r,l,b)},o.c))}),e(function(r,n,t){var 
e=f(xu,n,t.cD);if(1===e.$)return ke(t);var u=e.a;return 
f(Ne,function(){return ke(t)},f(Ne,function(n){return 
Ve(f(Te,function(t){return 
f(Re,r,t(n))},u))},mo))}),0,t(function(r,n){return 
f(_o,n.a,f(ze,r,n.b))}));var Xo=Kr("Time"),Yo=t(function(r,n){return 
Xo(f(_o,r,n))}),Qo=function(r){return 
r},Wo=Zr(L),rf=e(function(r,n,t){return{$:7,a:r,b:n,c:t}}),nf=function(r){return{$:4,a:r}},tf=t(function(r,n){return{$:10,a:r,b:n}}),ef=t(function(r,n){return 
f(ot,{X:f(ju,0,f(cu,gt(1),(t=f(Te,function(r){return 
r.X},r),t.b?ht(s(Lt,ne,t.a,t.b)):$t))),aF:n},r);var 
t}),uf=function(r){return A(r,{V:""})},af=function(r){return 
s(Lt,t(function(r,n){return 
r.$?n:f(ot,r.a,n)}),L,r)},cf=t(function(r,n){return 
s(De,t(function(n,t){return 
r(n)?f(ot,n,t):t}),L,n)}),of=u(function(r,n,t,e){var u=n;return 
f(Fu,r,{cV:Oi,dm:"GET",dF:N([f(Po,"from",e),f(Po,"dir",t?"f":"b")]),dG:N(["rooms",u,"messages"]),dM:Bi})}),ff=t(function(r,n){return 
l(of,r,n.F,1,n.bK)}),sf=t(function(r,n){return 
l(of,r,n.F,0,n.bm)}),lf=Jr,bf=e(function(r,n,t){return 
1===n.$?$t:1===t.$?$t:ht(f(r,n.a,t.a))}),df=t(function(r,n){var 
t=n.b;return y(r(n.a),t)}),vf=t(function(r,n){var t;return 
A(r,{bK:(t=n.bK,1===t.$?r.bK:t.a),C:Mi(E(r.C,n.bA))})}),pf=t(function(r,n){var 
t;return 
A(r,{C:Mi(E(r.C,n.bA)),bm:(t=n.bK,1===t.$?n.bm:t.a)})}),mf=e(function(r,n,t){return 
f(1===n?vf:pf,r,t)}),gf=function(r){return 
f(Fn,"application/json",f(qt,0,r))},hf=t(function(r,n){return 
n.$?$t:r(n.a)}),$f=function(r){return{$:4,a:r}},wf=function(r){return{$:7,a:r}},yf=function(r){return{$:5,a:r}},qf={$:0},xf={$:8},Af=t(function(r,n){return{$:4,a:r,b:n}}),kf=function(r){return{$:2,a:r}},Ef=function(r){return{$:0,a:r}},Lf=function(r){return{$:1,a:r}},Df=t(function(r,n){return{$:3,a:r,b:n}}),Tf=e(function(r,n,t){return{$:0,a:r,b:n,c:t}}),Nf=e(function(r,n,t){return{$:2,a:r,b:n,c:t}}),Sf=function(r){return{$:2,a:r}},Cf=e(function(r,n,t){return{$:1,a:r,b:n,c:t}}),Vf=t(function(r,n){return{$:0,a:r,b:n}}),Rf=function(r){return{$:1,a:r}},Uf=t(function(r,n){return{$:2,a:r,b:n}}),Bf=function(r){return{$:5,a:r}},Of=function(r){return{$:1,a:r}},Pf=function(r){return{$:2,a:r}},jf=t(function(r,n){return{$:6,a:r,b:n}}),If={$:8},_f={$:7},Gf=t(function(r,n){return{R:f(ot,n,r.R),H:r.H}}),zf={$:10},Hf=function(r){switch(r){case" 
":case"\t":return!0;default:return!1}},Mf=f(Xa,"\r",Ja("a carriage 
return")),Ff=f(Xa,"\n",Ja("a 
newline")),Kf=yc(N([nc(Ff),f(Na,nc(Mf),yc(N([nc(Ff),_a(0)])))])),Zf=f(uc,function(){return 
zf},f(Na,va(Ea(Hf)),Kf)),Jf=function(r){return{$:11,a:r}},Xf=f(Xa," 
",Ja("a 
space")),Yf=nc,Qf=N([Yf(f(Xa,">",Ja(">"))),f(Na,va(Yf(Xf)),yc(N([Yf(f(Xa,">",Ja(" 
>"))),Yf(f(Xa," >",Ja("  >"))),Yf(f(Xa,"  >",Ja("   
>")))])))]),Wf=function(r){switch(r){case"\n":case"\r":return!0;default:return!1}},rs=Ea(f(ze,Oa,Wf)),ns=Yi(Ja("the 
end of the 
input")),ts=yc(N([Kf,ns])),es=f(Qa,f(Na,f(Na,_a(Jf),yc(Qf)),yc(N([Yf(Xf),_a(0)]))),f(Na,Ra(rs),ts)),us=t(function(r,n){return{$:0,a:r,b:n}}),as=function(r){return{$:6,a:r}},cs=function(r){return{$:8,a:r}},is=t(function(r,n){return{$:0,a:r,b:n}}),os=t(function(r,n){return{$:0,a:r,b:n}}),fs=e(function(r,n,t){var 
e=y(n,t);return""===e.a?t:""===e.b?n:E(n,E(r,t))}),ss=t(function(r,n){return 
r+"\n"+n}),ls=function(r){return nc(f(Xa,r,Ja(r)))},bs=function(r){var 
n=r.a,e=r.b,u=f(ju,pc(e),f(cu,function(r){return 
pc(f(ot,r,e))},n)),a=f(ju,mc(y($t,e)),f(cu,function(r){return 
mc(y($t,f(ot,r,e)))},n)),c=function(r){return mc(y(ht(function(r){return 
E(f(ju,"",n),r)}(r)),e))};return yc(N([f(uc,function(){return 
u},ls("|\n")),f(uc,function(){return u},ls("\n")),f(uc,function(){return 
u},Yi(Ja("end"))),f(Na,va(_a(c("|"))),ls("\\\\|")),f(Na,va(_a(c("\\"))),ls("\\\\")),f(Na,va(_a(c("|"))),ls("\\|")),f(Na,va(_a(a)),ls("|")),f(Va,t(function(r){return 
c(r)}),f(xa,Da(!0),Pa("No character found")))]))},ds=function(r){return 
r.trim()},vs=f(uc,f(Lt,t(function(r,n){return 
f(ot,ds(r),n)}),L),f(vc,y($t,L),bs)),ps=f(Qa,f(Na,_a(pe),yc(N([ls("|"),_a(0)]))),vs),ms=t(function(r,n){var 
e=r.b,u=f(yi,ps,n);return u.$?lt("Unable to parse previous line as a table 
header"):function(r){return p(Dt(r),Dt(e))?pt(function(r){return 
s(Tt,t(function(r,n){return{at:n,Y:r}}),r,e)}(r)):lt("Tables must have the 
same number of header columns ("+xt(Dt(r))+") as delimiter columns 
("+xt(Dt(e))+")")}(u.a)}),gs=t(function(r,n){return{R:r.R,H:function(){var 
t=y(n,r.H);r:for(;t.b.b;)switch(t.b.a.$){case 5:if(5===t.a.$){var 
e=t.b,u=e.b;return f(ot,yf({cV:f(ss,e.a.a.cV,t.a.a.cV),di:$t}),u)}break 
r;case 6:if(6===t.a.$){var a=t.b;return 
u=a.b,f(ot,as(f(ss,a.a.a,t.a.a)),u)}break r;case 11:switch(t.a.$){case 
1:var c=t.b;return u=c.b,f(ot,Jf(s(fs,"\n",c.a.a,t.a.a)),u);case 11:var 
i=t.b;return u=i.b,f(ot,Jf(f(ss,i.a.a,t.a.a)),u);default:break r}case 
1:switch(t.a.$){case 1:var o=t.b;return 
u=o.b,f(ot,Rf(s(fs,"\n",o.a.a,t.a.a)),u);case 12:if(t.a.a){var 
l=t.b;return u=l.b,f(ot,f(us,2,l.a.a),u)}var b=t.b;u=b.b;return 
f(ot,f(us,1,b.a.a),u);case 9:var 
d=t.a.a,v=d.a,p=t.b,m=p.a.a,g=(u=p.b,f(ms,f(os,v,d.b),m));return 
f(ot,g.$?Rf(s(fs,"\n",m,v.cj)):cs(f(is,g.a,L)),u);default:break r}case 
8:if(8===t.a.$)return u=t.b.b,f(ot,cs(t.a.a),u);break r;default:break 
r}return f(ot,n,r.H)}()}}),hs=t(function(r,n){return 
n.b?s(De,ot,n,r):r}),$s=function(r){return 
s(De,hs,L,r)},ws=t(function(r,n){return 
$s(f(Te,r,n))}),ys=function(r){return"Problem at row 
"+xt(r.dP)+"\n"+function(r){switch(r.$){case 0:return"Expecting "+r.a;case 
1:return"Expecting int";case 2:return"Expecting hex";case 
3:return"Expecting octal";case 4:return"Expecting binary";case 
5:return"Expecting float";case 6:return"Expecting number";case 
7:return"Expecting variable";case 8:return"Expecting symbol "+r.a;case 
9:return"Expecting keyword "+r.a;case 10:return"Expecting keyword 
end";case 11:return"Unexpected char";case 12:return r.a;default:return"Bad 
repeat"}}(r.dH)},qs=function(r){return{$:3,a:r}},xs=e(function(r,n,t){return{$:0,a:r,b:n,c:t}}),As=Ja("at 
least 1 tag name character"),ks=function(r){switch(r){case" 
":case"\r":case"\n":case"\t":case"/":case"<":case">":case'"':case"'":case"=":return!1;default:return!0}},Es=f(Va,t(function(r){return 
za(r)}),f(Na,f(xa,ks,As),Ea(ks))),Ls=Es,Ds=function(r){return{$:8,a:r}},Ts=function(r){return 
nc(f(Xa,r,Ds(r)))},Ns=Ec(N([y("amp","&"),y("lt","<"),y("gt",">"),y("apos","'"),y("quot",'"')])),Ss=t(function(r,n){return 
n.$?lt(r):pt(n.a)}),Cs=function(r){var n=function(r){var 
n,t=f(xe,"#x",n=r)?f(Bc,Pa,f(Uc,Cc,jc(f($e,2,n)))):f(xe,"#",n)?f(Ss,Pa("Invalid 
escaped character: "+n),f(cu,Cc,qe(f($e,1,n)))):f(Ss,Pa('No entity named 
"&'+n+';" found.'),f(xu,n,Ns));return 
t.$?ja(t.a):_a(t.a)},t=function(n){return!p(n,r)&&";"!==n};return 
f(Qa,f(Na,_a(pe),Ts("&")),f(Na,f(ba,n,Ra(f(Na,f(xa,t,Ja("an entity 
character")),Ea(t)))),Ts(";")))},Vs=e(function(r,n,t){return 
f(ba,function(n){return yc(N([f(uc,function(r){return 
mc(E(t,E(n,Sc(r))))},Cs(r)),_a(pc(E(t,n)))]))},Ra(Ea(n)))}),Rs=function(r){return 
f(vc,"",f(Vs,r,function(n){return!p(n,r)&&"&"!==n}))},Us=yc(N([f(Qa,f(Na,_a(pe),Ts('"')),f(Na,Rs('"'),Ts('"'))),f(Qa,f(Na,_a(pe),Ts("'")),f(Na,Rs("'"),Ts("'")))])),Bs=t(function(r,n){return 
ht(n.$?r:n.a)}),Os=function(r){switch(r){case" 
":case"\r":case"\n":case"\t":return!0;default:return!1}},Ps=Ea(Os),js=function(r){return 
yc(N([f(Qa,f(Qa,_a(t(function(n,t){return 
mc(s(Uu,za(n),Bs(t),r))})),f(Na,f(Na,f(Na,Ls,Ps),Ts("=")),Ps)),f(Na,Us,Ps)),_a(pc(r))]))},Is=f(uc,f(Mo,e(function(r,n,t){return 
f(ot,{ai:r,aP:n},t)}),L),f(vc,$u,js)),_s=function(r){return 
function(n){var 
t=b(Wn,r,n.b,n.dP,n.bB,n.cy),e=t.a,u=t.b,a=t.c,c=0>e?ge(n.cy):e;return 
s(la,0>h(n.b,c),0,{bB:a,d:n.d,e:n.e,b:c,dP:u,cy:n.cy})}},Gs=f(Qa,f(Na,_a(pe),Ts("<![CDATA[")),f(Na,Ra(_s("]]>")),Ts("]]>"))),zs=t(function(r,n){return 
f(uc,function(r){return r(n)},yc(r))}),Hs=function(r){var 
n=f(ba,function(n){return p(r,n)?_a(0):ja(Pa("tag name mismatch: "+r+" and 
"+n))},Es);return 
f(Na,f(Na,f(Na,f(Na,Ts("</"),Ps),n),Ps),Ts(">"))},Ms=function(r){return 
f(Xa,r,Ja(r))},Fs=f(Qa,f(Na,_a(function(r){return{$:2,a:r}}),nc(Ms("\x3c!--"))),f(Na,Ra(_s("--\x3e")),nc(Ms("--\x3e")))),Ks=t(function(r,n){return{$:5,a:r,b:n}}),Zs=Ja("at 
least 1 uppercase 
character"),Js=Ra(f(Na,f(xa,Ut,Zs),Ea(Ut))),Xs=f(Na,f(xa,Os,Ja("at least 
one 
whitespace")),Ea(Os)),Ys=f(Qa,f(Qa,f(Na,_a(Ks),Ts("<!")),f(Na,Js,Xs)),f(Na,Ra(_s(">")),Ts(">"))),Qs=f(Qa,f(Na,_a(function(r){return{$:4,a:r}}),Ts("<?")),f(Na,Ra(_s("?>")),Ts("?>"))),Ws=function(r){switch(r){case"<":case"&":return!1;default:return!0}},rl=N([f(uc,function(){return 
mc(0)},f(Na,f(xa,Ws,Ja("is not & or <")),Ea(Ws))),f(uc,function(){return 
mc(0)},Cs("<")),_a(pc(0))]),nl=Ra(f(vc,0,function(){return 
yc(rl)})),tl=function(r){return f(vc,L,zs(el(r)))},el=function(r){return 
N([f(uc,t(function(r,n){return pc(jt(n))}),Hs(r)),f(ba,function(n){return 
we(n)?f(uc,t(function(r,n){return pc(jt(n))}),Hs(r)):_a(function(r){return 
mc(f(ot,{$:1,a:n},r))})},nl),f(uc,t(function(r,n){return 
mc(f(ot,r,n))}),al())])},ul=function(r){return 
f(Qa,f(Qa,f(Na,_a(xs(r)),Ps),f(Na,Is,Ps)),yc(N([f(uc,function(){return 
L},Ts("/>")),f(Qa,f(Na,_a(pe),Ts(">")),tl(r))])))};function al(){return 
yc(N([f(uc,qs,Gs),Qs,Fs,Ys,cl()]))}function cl(){return 
f(Qa,f(Na,_a(pe),Ts("<")),f(ba,ul,Es))}var il=al();al=function(){return 
il};var ol=cl();cl=function(){return ol};var 
fl,sl,ll,bl,dl=f(Xa,"\t",Ja("a 
tab")),vl=yc(N([Yf(dl),f(Na,va(Yf(Xf)),yc(N([Yf(f(Xa,"   
",Ds("Indentation"))),Yf(f(Xa," \t",Ds("Indentation"))),Yf(f(Xa,"  
\t",Ds("Indentation")))])))])),pl=f(Qa,f(Na,_a(as),vl),f(Na,Ra(rs),ts)),ml=t(function(r,n){return 
y(n.a,r(n.b))}),gl=f(Va,t(function(r){return 
Rf(r)}),rs),hl=f(Na,gl,ts),$l=t(function(r,n){return{$:4,a:r,b:n}}),wl=function(r){return 
f(Na,f(xa,r,Pa("Expected one or more 
character")),Ea(r))},yl=f(Xa,")",Ja("a `)`")),ql=f(Xa,".",Ja("a 
`.`")),xl=yc(N([f(Qa,f(Na,_a(pe),wl(Hf)),f(Na,Ra(rs),ts)),f(Na,_a(""),ts)])),Al=f(Va,t(function(r){return 
f(ju,0,qe(r))}),wl(Ot)),kl=t(function(r,n){return 
yc(N([f(uc,function(r){return 
mc(f(ot,r,n))},r),_a(pc(jt(n)))]))}),El=e(function(r,n,t){return 
f(uc,function(n){return y(r,f(ot,t,n))},f(vc,L,kl(function(r){return 
f(Qa,f(Na,_a(pe),va(f(Na,Al,Yf(r)))),xl)}(n))))}),Ll=f(ba,function(r){return 
r>999999999?ja(Pa("Starting numbers must be nine digits or 
less.")):_a(r)},Al),Dl=function(r){return 1===r?_a(r):ja(Pa("Lists inside 
a paragraph or after a paragraph without a blank line must start with 
1"))},Tl=function(r){return f(uc,function(r){return 
f($l,r.a,f(Te,pe,r.b))},function(r){return 
f(ba,pe,f(Qa,f(Qa,f(Qa,_a(El),va(r?f(ba,Dl,Ll):Ll)),f(Na,va(yc(N([f(Na,_a(ql),Yf(ql)),f(Na,_a(yl),Yf(yl))]))),wl(Hf))),f(Na,Ra(rs),ts)))}(r))},Nl=t(function(r,n){return{$:6,a:r,b:n}}),Sl={$:1},Cl=e(function(r,n,t){return{$:4,a:r,b:n,c:t}}),Vl=e(function(r,n,t){return{$:3,a:r,b:n,c:t}}),Rl=function(r){return{$:0,a:r}},Ul=function(r){var 
n=r,t=n.g;switch(t.$){case 0:return Rl(n.d5);case 1:return Sl;case 
2:return{$:2,a:n.d5};case 3:var e=t.a;return 
s(Vl,e.b,$t,N([Rl(e.a)]));case 4:var u=t.a;return 
s(Vl,u.a,u.b,Bl(n.j));case 5:var a=t.a;return s(Cl,a.a,a.b,Bl(n.j));case 
6:return{$:5,a:t.a};case 7:return 
f(Nl,t.a,Bl(n.j));default:return{$:7,a:Bl(n.j)}}},Bl=function(r){return 
f(Te,Ul,r)},Ol=t(function(r,n){return{bK:n.bK-r.l,j:n.j,bm:n.bm-r.l,d5:n.d5,t:n.t-r.l,l:n.l-r.l,g:n.g}}),Pl=t(function(r,n){return{bK:r.bK,j:f(ot,f(Ol,r,n),r.j),bm:r.bm,d5:r.d5,t:r.t,l:r.l,g:r.g}}),jl=function(r){var 
n=r;return{bK:n.bK,j:Il(n.j),bm:n.bm,d5:n.d5,t:n.t,l:n.l,g:n.g}},Il=function(r){var 
n=f(Hi,function(r){return r.bm},r);return 
n.b?s(_l,n.b,n.a,L):L},_l=e(function(r,n,t){for(;;){var e=n;if(!r.b)return 
f(ot,jl(e),t);var 
u=r.a,a=r.b;1>h(e.bK,u.bm)?(r=a,n=u,t=f(ot,jl(e),t)):0>h(e.bm,u.bm)&&h(e.bK,u.bK)>0?(r=a,n=f(Pl,e,u),t=t):(r=a,n=e,t=t)}}),Gl={$:0},zl=u(function(r,n,t,e){return{a:n,ah:r,dy:t,bo:e}}),Hl=nt,Ml=function(r){return 
f(Hl,{cX:!1,dp:!1},r)},Fl=/.^/,Kl=f(ju,Fl,Ml("&#([0-9]{1,8});")),Zl=ut(1/0),Jl=I,Xl=function(r){return 
function(r){return 
9===r||10===r||13===r||133===r||r>=32&&126>=r||r>=160&&55295>=r||r>=57344&&64975>=r||r>=65008&&65533>=r||r>=65536&&1114109>=r}(r)&&!function(r){var 
n=f(Jl,16,r),t=f(Jl,131070,r);return!(131070>r||(0>t||t>15)&&(65536>t||t>65551)||14!==n&&15!==n)}(r)?Sc(Cc(r)):Sc(Cc(65533))},Yl=f(Zl,Kl,function(r){var 
n=r.bo;if(n.b&&!n.a.$){var t=qe(n.a.a);return t.$?r.ah:Xl(t.a)}return 
r.ah}),Ql=f(ju,Fl,Ml("&([0-9a-zA-Z]+);")),Wl=Ec(N([y("quot",34),y("amp",38),y("apos",39),y("lt",60),y("gt",62),y("nbsp",160),y("iexcl",161),y("cent",162),y("pound",163),y("curren",164),y("yen",165),y("brvbar",166),y("sect",167),y("uml",168),y("copy",169),y("ordf",170),y("laquo",171),y("not",172),y("shy",173),y("reg",174),y("macr",175),y("deg",176),y("plusmn",177),y("sup2",178),y("sup3",179),y("acute",180),y("micro",181),y("para",182),y("middot",183),y("cedil",184),y("sup1",185),y("ordm",186),y("raquo",187),y("frac14",188),y("frac12",189),y("frac34",190),y("iquest",191),y("Agrave",192),y("Aacute",193),y("Acirc",194),y("Atilde",195),y("Auml",196),y("Aring",197),y("AElig",198),y("Ccedil",199),y("Egrave",200),y("Eacute",201),y("Ecirc",202),y("Euml",203),y("Igrave",204),y("Iacute",205),y("Icirc",206),y("Iuml",207),y("ETH",208),y("Ntilde",209),y("Ograve",210),y("Oacute",211),y("Ocirc",212),y("Otilde",213),y("Ouml",214),y("times",215),y("Oslash",216),y("Ugrave",217),y("Uacute",218),y("Ucirc",219),y("Uuml",220),y("Yacute",221),y("THORN",222),y("szlig",223),y("agrave",224),y("aacute",225),y("acirc",226),y("atilde",227),y("auml",228),y("aring",229),y("aelig",230),y("ccedil",231),y("egrave",232),y("eacute",233),y("ecirc",234),y("euml",235),y("igrave",236),y("iacute",237),y("icirc",238),y("iuml",239),y("eth",240),y("ntilde",241),y("ograve",242),y("oacute",243),y("ocirc",244),y("otilde",245),y("ouml",246),y("divide",247),y("oslash",248),y("ugrave",249),y("uacute",250),y("ucirc",251),y("uuml",252),y("yacute",253),y("thorn",254),y("yuml",255),y("OElig",338),y("oelig",339),y("Scaron",352),y("scaron",353),y("Yuml",376),y("fnof",402),y("circ",710),y("tilde",732),y("Alpha",913),y("Beta",914),y("Gamma",915),y("Delta",916),y("Epsilon",917),y("Zeta",918),y("Eta",919),y("Theta",920),y("Iota",921),y("Kappa",922),y("Lambda",923),y("Mu",924),y("Nu",925),y("Xi",926),y("Omicron",927),y("Pi",928),y("Rho",929),y("Sigma",931),y("Tau",932),y("Upsilon",933),y("Phi",934),y("Chi",935),y("Psi",936),y("Omega",937),y("alpha",945),y("beta",946),y("gamma",947),y("delta",948),y("epsilon",949),y("zeta",950),y("eta",951),y("theta",952),y("iota",953),y("kappa",954),y("lambda",955),y("mu",956),y("nu",957),y("xi",958),y("omicron",959),y("pi",960),y("rho",961),y("sigmaf",962),y("sigma",963),y("tau",964),y("upsilon",965),y("phi",966),y("chi",967),y("psi",968),y("omega",969),y("thetasym",977),y("upsih",978),y("piv",982),y("ensp",8194),y("emsp",8195),y("thinsp",8201),y("zwnj",8204),y("zwj",8205),y("lrm",8206),y("rlm",8207),y("ndash",8211),y("mdash",8212),y("lsquo",8216),y("rsquo",8217),y("sbquo",8218),y("ldquo",8220),y("rdquo",8221),y("bdquo",8222),y("dagger",8224),y("Dagger",8225),y("bull",8226),y("hellip",8230),y("permil",8240),y("prime",8242),y("Prime",8243),y("lsaquo",8249),y("rsaquo",8250),y("oline",8254),y("frasl",8260),y("euro",8364),y("image",8465),y("weierp",8472),y("real",8476),y("trade",8482),y("alefsym",8501),y("larr",8592),y("uarr",8593),y("rarr",8594),y("darr",8595),y("harr",8596),y("crarr",8629),y("lArr",8656),y("uArr",8657),y("rArr",8658),y("dArr",8659),y("hArr",8660),y("forall",8704),y("part",8706),y("exist",8707),y("empty",8709),y("nabla",8711),y("isin",8712),y("notin",8713),y("ni",8715),y("prod",8719),y("sum",8721),y("minus",8722),y("lowast",8727),y("radic",8730),y("prop",8733),y("infin",8734),y("ang",8736),y("and",8743),y("or",8744),y("cap",8745),y("cup",8746),y("int",8747),y("there4",8756),y("sim",8764),y("cong",8773),y("asymp",8776),y("ne",8800),y("equiv",8801),y("le",8804),y("ge",8805),y("sub",8834),y("sup",8835),y("nsub",8836),y("sube",8838),y("supe",8839),y("oplus",8853),y("otimes",8855),y("perp",8869),y("sdot",8901),y("lceil",8968),y("rceil",8969),y("lfloor",8970),y("rfloor",8971),y("lang",9001),y("rang",9002),y("loz",9674),y("spades",9824),y("clubs",9827),y("hearts",9829),y("diams",9830)])),rb=f(Zl,Ql,function(r){var 
n=r.bo;if(n.b&&!n.a.$){var t=f(xu,n.a.a,Wl);return 
t.$?r.ah:Sc(Cc(t.a))}return 
r.ah}),nb=f(ju,Fl,Ml("(\\\\+)([!\"#$%&\\'()*+,./:;<=>?@[\\\\\\]^_`{|}~-])")),tb=e(function(r,n,t){return 
r>0?s(tb,r>>1,E(n,n),1&r?E(t,n):t):t}),eb=t(function(r,n){return 
s(tb,r,n,"")}),ub=f(Zl,nb,function(r){var 
n=r.bo;if(n.b&&!n.a.$&&n.b.b&&!n.b.a.$){var t=n.b.a.a;return 
E(f(eb,ge(n.a.a)/2|0,"\\"),t)}return 
r.ah}),ab=f(ju,Fl,Ml("&#[Xx]([0-9a-fA-F]{1,8});")),cb=F,ib=f(Zl,ab,function(r){var 
n,e=r.bo;return e.b&&!e.a.$?Xl((n=e.a.a,s(cb,t(function(r,n){return 
16*n+f(Jl,39,Vt(r))-9}),0,za(n)))):r.ah}),ob=function(r){var 
n=ub(r);return 
f(me,"&",n)?ib(Yl(rb(n))):n},fb=function(r){return{bK:0,j:L,bm:0,d5:ob(r),t:0,l:0,g:Gl}},sb=e(function(r,n,t){var 
e=n,u={bK:e.bK,j:s(lb,e.d5,L,e.j),bm:e.bm,d5:e.d5,t:e.t,l:e.l,g:e.g};if(t.b){var 
a=t.a;return 
a.g.$?p(e.bK,a.bm)?f(ot,u,t):0>h(e.bK,a.bm)?f(ot,u,f(ot,fb(s(he,e.bK,a.bm,r)),t)):t:f(ot,u,t)}var 
c=f($e,e.bK,r);return 
we(c)?N([u]):N([u,fb(c)])}),lb=e(function(r,n,t){for(;;){if(!t.b){if(n.b){var 
e=n.a;return e.bm>0?f(ot,fb(f(ye,e.bm,r)),n):n}return 
we(r)?L:N([fb(r)])}var 
u=t.b,a=r,c=s(sb,r,t.a,n);r=a,n=c,t=u}}),bb=f(ju,Fl,Ml("(\\\\*)(\\<)")),db=e(function(r,n,t){var 
e=r(n);return e.$?t:f(ot,e.a,t)}),vb=t(function(r,n){return 
s(De,db(r),L,n)}),pb=et(1/0),mb={$:4},gb=function(r){return!f(Jl,2,r)},hb=function(r){var 
n=r.bo;if(n.b&&n.b.b&&!n.b.a.$){var t=f(ju,0,f(cu,ge,n.a));return 
gb(t)?ht({a:r.a+t,a3:1,c:mb}):$t}return 
$t},$b=f(ju,Fl,Ml("(\\\\*)(\\>)")),wb=function(r){return{$:5,a:r}},yb=function(r){var 
n=r.bo;if(n.b&&n.b.b&&!n.b.a.$){var t=f(ju,0,f(cu,ge,n.a));return 
ht({a:r.a+t,a3:1,c:gb(t)?wb(1):wb(0)})}return 
$t},qb=f(ju,Fl,Ml("(\\\\*)([^*])?(\\*+)([^*])?")),xb=t(function(r,n){return{$:7,a:r,b:n}}),Ab=f(cb,t(function(r,n){return 
n||function(r){switch(r){case"!":case'"':case"#":case"%":case"&":case"'":case"(":case")":case"*":case",":case"-":case".":case"/":case":":case";":case"?":case"@":case"[":case"]":case"_":case"{":case"}":case"~":return!0;default:return!1}}(r)}),!1),kb=f(cb,t(function(r,n){return 
n||function(r){switch(r){case" 
":case"\f":case"\n":case"\r":case"\t":case"\v":case" 
":case"\u2028":case"\u2029":return!0;default:return!1}}(r)}),!1),Eb=function(r){if(r.$)return 
0;var n=r.a;return we(n)||kb(n)?0:Ab(n)?1:2},Lb=e(function(r,n,t){var 
e=t.bo;if(e.b&&e.b.b&&e.b.b.b&&!e.b.b.a.$&&e.b.b.b.b){var 
u=e.a,a=e.b,c=a.a,i=a.b,o=i.a.a,l=Eb(i.b.a),b=c.$?0:ge(c.a),d=t.a&&!b?ht(s(he,t.a-1,t.a,n)):c,v=u.$?0:ge(u.a),p=!gb(v)&&!b||!d.$&&"\\"===d.a,m=p?ge(o)-1:ge(o),g=p?1:Eb(d);return 
0>=m||"_"===r&&2===g&&2===l?$t:ht({a:t.a+v+b+(p?1:0),a3:m,c:f(xb,r,{aD:g,aJ:l})})}return 
$t}),Db=f(ju,Fl,Ml("(\\\\*)(\\`+)")),Tb=function(r){return{$:0,a:r}},Nb=function(r){var 
n=r.bo;if(n.b&&n.b.b&&!n.b.a.$){var 
t=n.b.a.a,e=f(ju,0,f(cu,ge,n.a));return 
ht({a:r.a+e,a3:ge(t),c:gb(e)?Tb(1):Tb(0)})}return 
$t},Sb=f(ju,Fl,Ml("(?:(\\\\+)|( {2,}))\\n")),Cb={$:9},Vb=function(r){var 
n=r.bo;r:for(;;){if(n.b){if(n.a.$){if(n.b.b&&!n.b.a.$)return 
ht({a:r.a,a3:ge(r.ah),c:Cb});break r}var t=ge(n.a.a);return 
gb(t)?$t:ht({a:r.a+t-1,a3:2,c:Cb})}break r}return $t},Rb=(Ml("(?:(\\\\+)|( 
*))\\n"),f(ju,Fl,Ml("(\\\\*)(\\])"))),Ub={$:3},Bb=function(r){var 
n=r.bo;if(n.b&&n.b.b&&!n.b.a.$){var t=f(ju,0,f(cu,ge,n.a));return 
gb(t)?ht({a:r.a+t,a3:1,c:Ub}):$t}return 
$t},Ob=f(ju,Fl,Ml("(\\\\*)(\\!)?(\\[)")),Pb={$:2},jb=function(r){return{$:1,a:r}},Ib=function(r){var 
n=r.bo;if(n.b&&n.b.b&&n.b.b.b&&!n.b.b.a.$){var 
t=n.b.a,e=f(ju,0,f(cu,ge,n.a)),u=!gb(e),a=u?r.a+e+1:r.a+e;return 
u?t.$?$t:ht({a:a,a3:1,c:jb(0)}):ht(t.$?{a:a,a3:1,c:jb(0)}:{a:a,a3:2,c:Pb})}return 
$t},_b=function(r){return{$:10,a:r}},Gb=function(r){var 
n=r.bo;if(n.b&&n.b.b&&!n.b.a.$){var 
t=n.b.a.a,e=f(ju,0,f(cu,ge,n.a)),u=gb(e)?y(ge(t),_b(1)):y(ge(t),_b(0));return 
ht({a:r.a+e,a3:u.a,c:u.b})}return 
$t},zb=f(ju,Fl,Ml("(\\\\*)(~{2,})([^~])?")),Hb=f(ju,Fl,Ml("(\\\\*)([^_])?(\\_+)([^_])?")),Mb=t(function(r,n){if(r.b){var 
t=r.a,e=r.b;if(n.b){var u=n.a,a=n.b;return 
0>h(t.a,u.a)?f(ot,t,f(Mb,e,n)):f(ot,u,f(Mb,r,a))}return r}return 
n}),Fb={$:2},Kb=function(r){return{$:6,a:r}},Zb=function(r){return{$:5,a:r}},Jb=function(r){return{$:4,a:r}},Xb={$:8},Yb=function(r){return{$:3,a:r}},Qb=tt,Wb=f(ju,Fl,Ml("%(?:3B|2C|2F|3F|3A|40|26|3D|2B|24|23|25)")),rd=f(Tc,zu,f(Zl,Wb,function(r){return 
f(ju,r.ah,function(r){try{return ht(decodeURIComponent(r))}catch(r){return 
$t}}(r.ah))})),nd=f(ju,Fl,Ml("^([A-Za-z][A-Za-z0-9.+\\-]{1,31}:[^<>\\x00-\\x20]*)$")),td=et,ed=function(r){return 
r.b?ht(r.a):$t},ud=f(ju,Fl,Ml("^\\[\\s*([^\\[\\]\\\\]*(?:\\\\.[^\\[\\]\\\\]*)*)\\s*\\]")),ad=function(r){return 
r},cd=f(Tc,ad,za),id=t(function(r,n){return 
y(rd(ob(r)),f(cu,ob,n))}),od=e(function(r,n,t){var 
e,u=we(e=f(ju,r.d5,f(ju,$t,f(hf,f(Tc,function(r){return 
r.bo},ed),t))))?r.d5:e,a=f(xu,cd(u),n);if(1===a.$)return $t;var 
c=a.a,i=c.a,o=c.b,s=5===r.g.$?Zb(f(id,i,o)):Jb(f(id,i,o)),l=t.$?0:ge(t.a.ah);return 
ht(A(r,{bK:r.bK+l,g:s}))}),fd=e(function(r,n,t){return 
s(od,n,t,ed(s(td,1,ud,r)))}),sd=f(ju,Fl,Ml("^\\(\\s*(?:<([^<>\\f\\v\\r\\n]*)>|([^ 
\\t\\f\\v\\r\\n\\(\\)\\\\]*(?:\\\\.[^ \\t\\f\\v\\r\\n\\(\\)\\\\]*)*))(?:[ 
\\t\\f\\v\\r\\n]+(?:'([^'\\\\]*(?:\\\\.[^'\\\\]*)*)'|"+'"([^"\\\\]*(?:\\\\.[^"\\\\]*)*)"|\\(([^\\)\\\\]*(?:\\\\.[^\\)\\\\]*)*)\\)))?\\s*\\)')),ld=function(r){return 
s(Lt,t(function(r,n){return n.$?r:ht(n.a)}),$t,r)},bd=t(function(r,n){var 
t,e=n.bo;if(e.b&&e.b.b&&e.b.b.b&&e.b.b.b.b&&e.b.b.b.b.b){var 
u=e.a,a=e.b,c=a.a,i=a.b,o=i.b,s=ld(N([i.a,o.a,o.b.a])),l=ld(N([u,c]));return 
ht((t=f(ju,"",l),A(r,{bK:r.bK+ge(n.ah),g:(5===r.g.$?Zb:Jb)(f(id,t,s))})))}return 
$t}),dd=e(function(r,n,t){var e=n,u=s(td,1,sd,r);if(u.b){var 
a=f(bd,e,u.a);return a.$?s(fd,r,e,t):ht(a.a)}return 
s(fd,r,e,t)}),vd=t(function(r,n){var t=r,e=cf(function(r){var n=r;return 
h(t.bK,n.bm)>0&&0>h(t.bK,n.bK)});return 
ei(n)||ei(e(n))?ht(f(ot,t,n)):$t}),pd=f(ju,Fl,Ml("^([a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~\\-]+@[a-zA-Z0-9](?:[a-zA-Z0-9\\-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9\\-]{0,61}[a-zA-Z0-9])?)*)$")),md=function(r){var 
n=r;return 
f(Qb,pd,n.d5)?pt(A(n,{g:Yb(y(n.d5,"mailto:"+rd(n.d5)))})):lt(n)},gd=e(function(r,n,t){for(;;){if(!t.b)return 
$t;var e=t.a,u=t.b;if(n(e))return 
ht(q(e,jt(r),u));r=f(ot,e,r),n=n,t=u}}),hd=t(function(r,n){return 
s(gd,L,r,n)}),$d=t(function(r,n){return{$:6,a:r,b:n}}),wd=function(r){return 
s(la,!1,r.b,r)},yd=t(function(r,n){var 
t=n,u=f(Qa,f(Qa,f(Qa,_a(e(function(r,n,t){return{bW:n,a3:t-r}})),wd),il),wd),a=f(yi,u,f($e,t.bm,r));return 
a.$?$t:ht({a:t.bm,a3:a.a.a3,c:f($d,1,a.a.bW)})}),qd=t(function(r,n){return 
n.$?r(n.a):n}),xd=t(function(){return!1}),Ad=t(function(r,n){var 
t=n.c;return!t.$&&p(t.a?n.a3:n.a3-1,r.a3)}),kd=function(r){switch(r.c.$){case 
1:case 2:return!0;default:return!1}},Ed=t(function(r,n){var 
t=n.c;if(7===t.$){var e=t.b,u=r.c;if(7===u.$){var 
a=u.b;return!(!p(t.a,u.a)||(p(e.aD,e.aJ)||p(a.aD,a.aJ))&&!f(Jl,3,r.a3+n.a3))}return!1}return!1}),Ld=t(function(r,n){var 
t,e,u=10===(e=n.c).$?y(!0,e.a?n.a3:n.a3-1):y(!1,0),a=u.a,c=u.b,i=10===(t=r.c).$?y(!0,t.a?r.a3:r.a3-1):y(!1,0);return 
i.a&&a&&p(i.b,c)}),Dd={$:1},Td=t(function(r,n){return{bK:r.a+r.a3,j:L,bm:r.a,d5:"",t:0,l:0,g:n}}),Nd=a(function(r,n,t,e,u){for(;;){if(!r.b)return 
t;var 
a=r.a,c=r.b;9!==a.c.$?(r=c,n=f(ot,a,n),t=t,e=e,u=u):(r=c,n=n,t=f(ot,f(Td,a,Dd),t),e=e,u=u)}}),Sd=t(function(r,n){var 
t=r;return f(cf,function(r){return 
h(r.a,t.bK)>-1},n)}),Cd=c(function(r,n,t,e,u,a){var 
c,i=a.c,o=f(qd,md,f(Qb,nd,(c=v(Gd,e,u,function(r){return 
r},Fb,a.a,r,L)).d5)?pt(A(c,{g:Yb(y(c.d5,rd(c.d5)))})):lt(c));if(1===o.$){if(1===n){var 
s=f(yd,u,o.a);return s.$?$t:ht(y(f(ot,s.a,i),t))}return $t}return 
ht(y(i,f(ot,o.a,t)))}),Vd=a(function(r,n,t,e,u){r:for(;;){if(!r.b)return 
b(Od,jt(n),L,t,e,u);var a=r.a,c=r.b,i=a.c;switch(i.$){case 0:var 
o=i.a,s=f(hd,Ad(a),n);if(s.$){r=c,n=f(ot,a,n),t=t,e=e,u=u;continue r}var 
l=b(Rd,a,t,e,u,s.a);r=c,n=h=l.a,t=$=l.b,e=e,u=u;continue r;case 
5:o=i.a;var v=function(r){return 
4===r.c.$},p=f(hd,v,n);if(p.$){r=c,n=f(cf,f(ze,Oa,v),n),t=t,e=e,u=u;continue 
r}var 
m=d(Cd,a,o,t,e,u,p.a);if(m.$){r=c,n=f(cf,f(ze,Oa,v),n),t=t,e=e,u=u;continue 
r}var g=m.a,h=g.a,$=g.b;r=c,n=f(cf,f(ze,Oa,v),h),t=$,e=e,u=u;continue 
r;default:r=c,n=f(ot,a,n),t=t,e=e,u=u;continue 
r}}}),Rd=a(function(r,n,t,e,u){var 
a,c=u.a,i=u.c,o=(a=c.c).$||a.a?c:A(c,{a:c.a+1,a3:c.a3-1});return 
y(i,f(ot,v(Gd,t,e,ad,Fb,o,r,L),n))}),Ud=a(function(r,n,t,e,u){r:for(;;){if(!r.b)return 
b(Id,jt(n),L,t,e,u);var 
a=r.a,c=r.b,i=a.c;if(7!==i.$)r=c,n=f(ot,a,n),t=t,e=e,u=u;else{var 
o=i.a,s=i.b.aD,l=i.b.aJ;if(p(s,l)){if(!l||"_"===o&&1!==l){r=c,n=n,t=t,e=e,u=u;continue 
r}var d=f(hd,Ed(a),n);if(d.$){r=c,n=f(ot,a,n),t=t,e=e,u=u;continue r}var 
v=b(Bd,e,u,a,c,d.a);r=v.a,n=v.c,t=f(ot,v.b,t),e=e,u=u;continue 
r}if(0>h(s,l)){r=c,n=f(ot,a,n),t=t,e=e,u=u;continue r}var 
m=f(hd,Ed(a),n);if(m.$){r=c,n=n,t=t,e=e,u=u;continue r}var 
g=b(Bd,e,u,a,c,m.a);r=g.a,n=g.c,t=f(ot,g.b,t),e=e,u=u}}}),Bd=a(function(r,n,t,e,u){var 
a=u.a,c=u.b,i=u.c,o=a.a3-t.a3,s=o?o>0?{ax:t,aj:A(a,{a:a.a+o,a3:t.a3}),aI:f(ot,A(a,{a3:o}),i),aN:e}:{ax:A(t,{a3:a.a3}),aj:a,aI:i,aN:f(ot,A(t,{a:t.a+a.a3,a3:-o}),e)}:{ax:t,aj:a,aI:i,aN:e},l=v(Gd,r,n,function(r){return 
r},{$:7,a:s.aj.a3},s.aj,s.ax,jt(c));return 
q(s.aN,l,s.aI)}),Od=a(function(r,n,t,e,u){r:for(;;){if(!r.b)return 
b(Pd,jt(n),L,t,e,u);var 
a=r.a,c=r.b,i=a.c;if(6!==i.$)r=c,n=f(ot,a,n),t=t,e=e,u=u;else{var 
o=i.b;if(1===i.a){r=c,n=n,t=f(ot,f(Td,a,Kb(o)),t),e=e,u=u;continue r}var 
s=f(hd,xd(o),c);if(1===s.$){r=c,n=n,t=f(ot,f(Td,a,Kb(o)),t),e=e,u=u;continue 
r}var l=s.a,d=l.a,p=l.b,m=l.c,g=v(Gd,e,u,function(r){return 
r},Kb(o),a,d,p);r=m,n=n,t=f(ot,g,t),e=e,u=u}}}),Pd=a(function(r,n,t,e,u){r:for(;;){if(!r.b)return 
b(Ud,jt(n),L,t,e,u);var 
a=r.a,c=r.b;if(3!==a.c.$)r=c,n=f(ot,a,n),t=t,e=e,u=u;else{var 
i=f(hd,kd,n);if(i.$){r=c,n=n,t=t,e=e,u=u;continue r}var 
o=d(jd,a,c,t,e,u,i.a);if(o.$){r=c,n=n,t=t,e=e,u=u;continue r}var 
s=o.a;r=s.a,n=s.c,t=s.b,e=e,u=u}}}),jd=c(function(r,n,t,e,u,a){var 
c=a.a,i=a.b,o=a.c,l=q(n,t,E(i,o)),b=f($e,r.a+1,u),d=function(r){return 
1===r.c.$?A(r,{c:jb(1)}):r},p=function(n){return 
v(Gd,e,u,function(r){return 
r},n?Jb(y("",$t)):Zb(y("",$t)),c,r,jt(i))},m=c.c;switch(m.$){case 2:var 
g=p(!1),h=s(dd,b,g,e);if(1===h.$)return ht(l);var 
$=f(vd,k=h.a,t);if($.$)return ht(l);var w=$.a;return 
ht(q(f(Sd,k,n),w,o));case 1:if(m.a)return ht(l);g=p(!0);var 
x=s(dd,b,g,e);if(1===x.$)return ht(l);var k,L=f(vd,k=x.a,t);return 
L.$?ht(l):(w=L.a,ht(q(f(Sd,k,n),w,f(Te,d,o))));default:return 
$t}}),Id=a(function(r,n,t,e,u){r:for(;;){if(!r.b)return 
b(Nd,jt(n),L,t,e,u);var 
a=r.a,c=r.b;if(10!==a.c.$)r=c,n=f(ot,a,n),t=t,e=e,u=u;else{var 
i=f(hd,Ld(a),n);if(i.$){r=c,n=f(ot,a,n),t=t,e=e,u=u;continue r}var 
o=b(_d,a,t,e,u,i.a);r=c,n=o.a,t=o.b,e=e,u=u}}}),_d=a(function(r,n,t,e,u){var 
a,c=u.a,i=u.c,o=10!==(a=c.c).$||a.a?c:A(c,{a:c.a+1,a3:c.a3-1});return 
y(i,f(ot,v(Gd,t,e,ad,Xb,o,r,L),n))}),Gd=i(function(r,n,t,e,u,a,c){var 
i=u.a+u.a3,o=a.a,b=t(s(he,i,o,n)),d=u.a,v=a.a+a.a3,p={bK:v,j:L,bm:d,d5:b,t:o,l:i,g:e};return{bK:v,j:f(Te,function(r){return 
f(Ol,p,r)},l(zd,c,L,r,n)),bm:d,d5:b,t:o,l:i,g:e}}),zd=u(function(r,n,t,e){return 
b(Vd,r,L,n,t,e)}),Hd=t(function(r,n){var t=ds(n),e=function(r){return 
f(Mb,f(vb,yb,f(pb,$b,r)),f(Mb,f(vb,hb,f(pb,bb,r)),f(Mb,f(vb,Vb,f(pb,Sb,r)),f(Mb,f(vb,Bb,f(pb,Rb,r)),f(Mb,f(vb,Ib,f(pb,Ob,r)),f(Mb,f(vb,Gb,f(pb,zb,r)),f(Mb,(n=r,f(vb,f(Lb,"_",n),f(pb,Hb,n))),f(Mb,function(r){return 
f(vb,f(Lb,"*",r),f(pb,qb,r))}(r),f(vb,Nb,f(pb,Db,r))))))))));var 
n}(t);return Bl(s(lb,t,L,Il(l(zd,e,L,r,t))))}),Md=yc(N([nc(f(Xa," ",Ja(" 
"))),nc(f(Xa,">",Ja(">"))),f(Na,f(Na,f(xa,Bt,Ja("Alpha")),Ea(function(r){return 
Pt(r)||"-"===r})),yc(N([nc(f(Xa,":",Ja(":"))),nc(f(Xa,"@",Ja("@"))),nc(f(Xa,"\\",Ja("\\"))),nc(f(Xa,"+",Ja("+"))),nc(f(Xa,".",Ja(".")))])))])),Fd=va(f(Va,t(function(r){return 
Rf(r)}),f(Na,f(Na,f(Na,nc(f(Xa,"<",Ja("<"))),Md),rs),ts))),Kd=t(function(r,n){return{cV:n,di:r}}),Zd={av:"`",aC:0,aM:f(Xa,"`",Ja("a 
'`'"))},Jd=function(r){switch(r){case 1:return _a(0);case 2:return 
_a(1);case 3:return _a(2);case 4:return _a(3);default:return ja(Ja("Fenced 
code blocks should be indented no more than 3 
spaces"))}},Xd=e(function(r,n,t){for(;;){if(0>=n)return 
r;r=f(ot,t,r),n-=1,t=t}}),Yd=t(function(r,n){return 
s(Xd,L,r,n)}),Qd=t(function(r,n){var e=s(Lt,t(function(r,n){return 
f(Na,n,r)}),_a(0),f(Yd,r,nc(n.aM)));return f(Va,t(function(r){return 
y(n,ge(r))}),f(Na,e,Ea(Qt(n.av))))}),Wd=function(r){return 
s(la,!1,r.bB,r)},rv={av:"~",aC:1,aM:f(Xa,"~",Ja("a 
`~`"))},nv=nc(Xf),tv=yc(N([f(Na,f(Na,nv,yc(N([nv,_a(0)]))),yc(N([nv,_a(0)]))),_a(0)])),ev=f(Qa,f(Qa,f(Na,_a(t(function(r,n){return{aw:n.a,a0:r,a3:n.b}})),tv),f(ba,Jd,Wd)),yc(N([f(Qd,3,Zd),f(Qd,3,rv)]))),uv=Qt(" 
"),av=t(function(r,n){return 
f(Na,f(Na,f(Na,f(Na,_a(0),tv),f(Qd,r,n)),Ea(uv)),ts)}),cv=t(function(r,n){var 
e=f(Yd,r,n);if(e.b){var u=e.a,a=e.b;return s(Lt,t(function(r,n){return 
yc(N([f(Na,r,n),_a(0)]))}),yc(N([u,_a(0)])),a)}return 
_a(0)}),iv=function(r){return s(la,!1,r.cy,r)},ov=function(r){var 
n,u=r.a,a=r.b;return 
yc(N([f(Na,_a(pc(a)),Yi(Xi)),f(Va,t(function(r){return 
mc(y(u,E(a,r)))}),Kf),va(f(Na,_a(pc(a)),f(av,u.a3,u.aw))),f(Qa,f(Qa,f(Qa,_a(e(function(r,n,t){return 
mc(y(u,E(a,s(he,r,n,t))))})),(n=u.a0,f(Qa,f(Na,_a(pe),f(cv,n,nv)),f(Na,f(Na,wd,rs),ts)))),wd),iv)]))},fv=f(ba,function(r){return 
f(Qa,f(Qa,_a(Kd),f(Na,(n=r.aw,f(Va,t(function(r){var 
n=ds(r);return""===n?$t:ht(n)}),Ea(n.aC?f(ze,Oa,Wf):function(r){return"`"!==r&&!Wf(r)}))),ts)),function(r){return 
f(vc,y(r,""),ov)}(r));var n},ev),sv=t(function(r,n){return 
1>r?n:s(he,0,-r,n)}),lv=rr,bv=function(r){return 
f(lv,"#",r)?bv(f(sv,1,r)):r},dv=f(Xa,"#",Ja("a 
`#`")),vv=Ea(function(r){return" 
"===r||"\n"===r||"\r"===r}),pv=f(Qa,f(Qa,f(Na,f(Na,_a(us),f(ba,function(r){var 
n=ge(r);return 4>n?_a(n):ja(Ja("heading with < 4 spaces in 
front"))},Ra(vv))),Yf(dv)),f(ba,function(r){var n=ge(r)+1;return 
7>n?_a(n):ja(Ja("heading with < 7 
#'s"))},Ra(Ea(function(r){return"#"===r})))),yc(N([f(Na,_a(""),Yf(Ff)),f(Qa,f(Na,_a(pe),yc(N([Yf(Xf),Yf(dl)]))),f(Va,t(function(r){return 
n=ds(r),t=bv(n),f(lv," ",t)||we(t)?t.replace(/\s+$/,""):n;var 
n,t}),rs))]))),mv=f(Xa,">",Ja("a 
`>`")),gv=e(function(r,n,t){return{bB:n,d:t,dP:r}}),hv=t(function(r,n){return{bB:n.bB,d:r,e:n.e,b:n.b,dP:n.dP,cy:n.cy}}),$v=t(function(r,n){var 
t=n;return function(n){var e=t(f(hv,f(ot,s(gv,n.dP,n.bB,r),n.d),n));return 
e.$?e:s(la,e.a,e.b,f(hv,n.d,e.c))}}),wv=function(r){switch(r){case" 
":case"\n":case"\t":case"\v":case"\f":case"\r":return!0;default:return!1}},yv=f(Xa,"<",Ja("a 
`<`")),qv=f($v,"link 
destination",yc(N([f(Qa,f(Na,_a(zu),Yf(yv)),f(Na,Ra(Za(mv)),Yf(mv))),Ra(wl(f(ze,Oa,wv)))]))),xv=f(Xa,"]",Ja("a 
`]`")),Av=f(Xa,"[",Ja("a 
`[`")),kv=f(Qa,f(Na,_a(cd),Yf(Av)),f(Na,Ra(Za(xv)),Yf(f(Xa,"]:",Ja("]:"))))),Ev=f(Xa,'"',Ja("a 
double quote")),Lv=function(r){return f(me,"\n\n",r)?ja(Ja("no blank 
line")):_a(r)},Dv=f(Na,Ea(function(r){return!Wf(r)&&wv(r)}),ts),Tv=f(Na,f(xa,wv,Ja("Required 
whitespace")),Ea(wv)),Nv=f(Xa,"'",Ja("a single 
quote")),Sv=(fl=f(Qa,f(Na,_a(ht),Yf(Nv)),f(Na,f(Na,f(ba,Lv,Ra(Za(Nv))),Yf(Nv)),Dv)),sl=f(Qa,f(Na,_a(ht),Yf(Ev)),f(Na,f(Na,f(ba,Lv,Ra(Za(Ev))),Yf(Ev)),Dv)),f($v,"title",yc(N([va(f(Qa,f(Na,_a(pe),Tv),yc(N([sl,fl,_a($t)])))),f(Na,_a($t),Dv)])))),Cv=f($v,"link 
reference definition",f(Qa,f(Qa,f(Qa,f(Na,_a(e(function(r,n,t){return 
y(r,{bH:n,cG:t})})),tv),f(Na,f(Na,f(Na,kv,Ea(Hf)),yc(N([Kf,_a(0)]))),Ea(Hf))),qv),Sv)),Vv=Ea(Hf),Rv=function(r){var 
n=ls(Sc(r));return 
f(Na,f(Na,f(Na,f(Na,f(Na,f(Na,f(Na,_a(0),n),Vv),n),Vv),n),Ea(function(n){return 
p(n,r)||Hf(n)})),ts)},Uv=yc(N([Rv("-"),Rv("*"),Rv("_")])),Bv=yc(N([f(Qa,f(Na,f(Na,f(Na,_a(pe),nv),yc(N([nv,_a(0)]))),yc(N([nv,_a(0)]))),Uv),Uv])),Ov=t(function(r,n){return{$:12,a:r,b:n}}),Pv=f(Xa,"=",Ja("a 
`=`")),jv=f(Xa,"-",Ja("a `-`")),Iv=(ll=e(function(r,n,t){return 
f(Na,f(Na,_a(r),nc(n)),Ea(Qt(t)))}),f(Va,t(function(r,n){return 
f(Ov,n,r)}),f(Qa,f(Na,_a(pe),tv),f(Na,f(Na,yc(N([s(ll,0,Pv,"="),s(ll,1,jv,"-")])),Ea(Hf)),ts)))),_v=function(r){return{$:9,a:r}},Gv=Ea(Hf),zv=function(r){return 
yc(N([f(xa,r,Pa("Character not found")),_a(0)]))},Hv=function(r){return 
yc(N([va(f(uc,function(){return pc(r)},ls("|\n"))),f(uc,function(){return 
pc(r)},ls("\n")),f(uc,function(){return 
pc(r)},Yi(Ja("end"))),va(f(Na,f(Na,_a(pc(r)),ls("|")),Yi(Ja("end")))),f(Qa,f(Na,f(Na,_a(function(n){return 
mc(f(ot,n,r))}),(n=r,ei(n)?yc(N([ls("|"),_a(0)])):ls("|"))),Gv),f(Na,Ra(f(Na,f(Na,f(Na,_a(0),zv(function(r){return":"===r})),wl(function(r){return"-"===r})),zv(function(r){return":"===r}))),Gv))]));var 
n},Mv=function(r){var n=y(f(xe,":",r),f(lv,":",r));return 
n.a?ht(n.b?2:0):n.b?ht(1):$t},Fv=f(uc,_v,f(ba,function(r){var 
n=r.a.cJ,t=r.b;return ei(t)?ja(Ja("Must have at least one column in 
delimiter row.")):1!==Dt(t)||f(xe,"|",n)&&f(lv,"|",n)?_a(r):ja(Pa("Tables 
with a single column must have pipes at the start and end of the delimiter 
row to avoid ambiguity."))},f(Va,t(function(r,n){return 
f(os,{cj:r,cJ:ds(r)},f(Te,Mv,jt(n)))}),f(vc,L,Hv)))),Kv=t(function(r,n){return!f(ic,f(ze,Oa,r),n)}),Zv=e(function(r,n,t){r:for(;;){if(r>0){if(n.b){var 
e=n.a;r-=1,n=n.b,t=f(ot,e,t);continue r}return t}return 
t}}),Jv=t(function(r,n){return 
jt(s(Zv,r,n,L))}),Xv=e(function(r,n,t){if(n>0){var 
e=y(n,t);r:for(;;){n:for(;;){if(!e.b.b)return 
t;if(!e.b.b.b){if(1===e.a)break r;break n}switch(e.a){case 1:break r;case 
2:var u=e.b;return N([u.a,u.b.a]);case 3:if(e.b.b.b.b){var 
a=e.b,c=a.b;return N([a.a,c.a,c.b.a])}break 
n;default:if(e.b.b.b.b&&e.b.b.b.b.b){var 
i=e.b,o=i.b,l=o.b,b=l.b,d=b.b;return 
f(ot,i.a,f(ot,o.a,f(ot,l.a,f(ot,b.a,r>1e3?f(Jv,n-4,d):s(Xv,r+1,n-4,d)))))}break 
n}}return t}return N([e.b.a])}return L}),Yv=t(function(r,n){return 
s(Xv,0,r,n)}),Qv=t(function(r,n){var t=Dt(n);switch(f(qu,r,t)){case 
0:return f(Yv,r,n);case 1:return n;default:return 
E(n,f(Yd,r-t,""))}}),Wv=function(r){var n,t=r.a,e=r.b;return 
f(uc,function(r){return 
cs(f(is,t,E(e,N([r]))))},(n=Dt(t),f(ba,function(r){return 
ei(r)||f(Kv,we,r)?ja(Pa("A line must have at least one 
column")):_a(f(Qv,n,r))},ps)))},rp=function(r){return{$:3,a:r}},np=f(Xa,"*",Ja("a 
`*`")),tp=f(Xa,"+",Ja("a 
`+`")),ep=yc(N([f(Na,_a(jv),Yf(jv)),f(Na,_a(tp),Yf(tp)),f(Na,_a(np),Yf(np))])),up=function(r){return{$:1,a:r}},ap=t(function(r,n){return{$:0,a:r,b:n}}),cp=yc(N([f(Na,_a(1),Yf(f(Xa,"[x] 
",Ds("[x] ")))),f(Na,_a(1),Yf(f(Xa,"[X] ",Ds("[X] 
")))),f(Na,_a(0),Yf(f(Xa,"[ ] ",Ds("[ ] 
"))))])),ip=f(Qa,yc(N([f(Qa,_a(ap),f(Na,cp,Ea(Hf))),_a(up)])),f(Na,Ra(rs),ts)),op=yc(N([f(Qa,f(Na,_a(pe),va(wl(Hf))),ip),f(Na,_a(up("")),Kf)])),fp=e(function(r,n,t){return 
yc(N([f(uc,function(r){return 
mc(f(ot,r,t))},r),_a(pc(f(ot,n,jt(t))))]))}),sp=(bl=t(function(r,n){return 
f(vc,L,f(fp,function(r){return 
f(Qa,f(Na,_a(pe),va(Yf(r))),op)}(r),n))}),f(ba,pe,f(Qa,f(Qa,_a(bl),f(Na,va(ep),wl(Hf))),ip))),lp=f(uc,f(Tc,Te(function(r){return 
r.$?{cV:r.a,bq:$t}:{cV:r.b,bq:ht(1===r.a)}}),rp),sp),bp=t(function(r,n){switch(r.$){case 
0:var t=r.a,e=r.b,u=mp(r.c);if(u.$)return lt(u.a);var 
a=Ef(s(Tf,t,e,u.a));return pt(f(ot,a,n));case 1:var c=hp(r.a);return 
c.$?lt(Ja(f(At,"\n",f(Te,ys,c.a)))):pt(E(jt(c.a),n));case 2:return 
pt(f(ot,Ef(Lf(r.a)),n));case 3:return pt(f(ot,Ef($f(r.a)),n));case 
4:return pt(f(ot,Ef(Pf(r.a)),n));default:return 
pt(f(ot,Ef(f(Df,r.a,r.b)),n))}}),dp=t(function(r,n){var 
t=n,e=Ec(f(Te,ml(function(r){return y(r.bH,r.cG)}),r));return 
f(Te,vp,f(Hd,e,t))}),vp=function(r){switch(r.$){case 
0:return{$:7,a:r.a};case 1:return xf;case 2:return{$:6,a:r.a};case 
3:return s(Cf,r.a,r.b,f(Te,vp,n=r.c));case 4:return 
s(Nf,r.a,r.b,f(Te,vp,n=r.c));case 5:return{$:0,a:pp(r.a)};case 6:var 
n=r.b;switch(r.a){case 1:return{$:3,a:f(Te,vp,n)};case 2:default:return 
function(r){return{$:4,a:r}}(f(Te,vp,n))}default:return{$:5,a:f(Te,vp,n=r.a)}}},pp=function(r){switch(r.$){case 
1:return Lf("TODO this never happens, but use types to drop this 
case.");case 0:return s(Tf,r.a,r.b,f(ws,function(r){return 
1===r.$?Ep(r.a):N([Ef(pp(r))])},r.c));case 2:return Lf(r.a);case 3:return 
$f(r.a);case 4:return Pf(r.a);default:return 
f(Df,r.a,r.b)}},mp=function(r){return 
f(gp,r,L)},gp=t(function(r,n){for(;;){if(!r.b)return pt(jt(n));var 
t=r.b,e=f(bp,r.a,n);if(e.$)return lt(e.a);r=t,n=e.a}}),hp=function(r){var 
n=f(yi,f(Na,Sp(),ns),r);if(1===n.$)return lt(n.a);var t=$p(n.a);return 
1===t.$?f(yi,ja(t.a),""):pt(f(cf,function(r){return!(5===r.$&&!r.a.b)},t.a))},$p=function(r){return 
s(wp,r,r.H,L)},wp=e(function(r,n,t){r:for(;;){if(!n.b)return pt(t);var 
e=n.b,u=f(qp,r.R,n.a);switch(u.$){case 1:r=r,n=e,t=f(ot,u.a,t);continue 
r;case 0:r=r,n=e,t=t;continue r;default:return 
lt(u.a)}}}),yp=t(function(r,n){return f(Te,function(n){var t=n.at;return 
s(xp,r,function(r){return{at:t,Y:r}},n.Y)},n)}),qp=t(function(r,n){switch(n.$){case 
0:var t=n.b,e=function(r){switch(r){case 1:return pt(0);case 2:return 
pt(1);case 3:return pt(2);case 4:return pt(3);case 5:return pt(4);case 
6:return pt(5);default:return lt(Ja("A heading with 1 to 6 #'s, but found 
"+xt(r)))}}(n.a);return e.$?Sf(e.a):Of(f(Af,e.a,f(dp,r,t)));case 1:return 
Of(Bf(f(dp,r,t=n.a)));case 2:return Of(Ef(n.a));case 3:return 
Of({$:1,a:f(Te,function(n){var t,e=(t=n.bq).$?0:t.a?2:1;return 
f(Vf,e,s(xp,r,pe,n.cV))},n.a)});case 4:return 
t=n.b,Of(f(Uf,n.a,f(Te,f(xp,r,pe),t)));case 5:return Of(wf(n.a));case 
7:return Of(If);case 10:return qf;case 11:var 
u=n.a,a=f(yi,Sp(),u);if(a.$)return Sf(Pa(f(At,"\n",f(Te,ys,a.a))));var 
c=$p(a.a);return c.$?Sf(c.a):Of({$:3,a:c.a});case 6:return 
Of(wf({cV:n.a,di:$t}));case 8:var i=n.a,o=i.b;return 
Of(f(jf,f(yp,r,i.a),f(Ap,r,o)));case 9:return 
Of(Bf(f(dp,r,n.a.a.cj)));default:return 
Of(Bf(f(dp,r,n.b)))}}),xp=e(function(r,n,t){return 
n(f(dp,r,t))}),Ap=t(function(r,n){return f(Te,function(n){return 
f(Te,function(n){return s(xp,r,pe,n)},n)},n)}),kp=function(r){return 
yc(N([f(uc,function(){return pc(r)},ns),f(uc,function(n){return 
mc(f(Gf,r,n))},va(Cv)),f(uc,function(n){return 
mc(f(gs,r,n))},function(){var n=r.H;r:for(;n.b;)switch(n.a.$){case 
1:return Tp();case 8:var t=n.a.a;return yc(N([Dp(),Wv(t)]));default:break 
r}return Dp()}()),f(uc,function(n){return 
mc(f(gs,r,n))},hl)]))},Ep=function(r){return 
f(eu,L,hp(r))},Lp=function(r){switch(r.$){case 1:return _a(Rf(r.a));case 
0:var n=r.a,t=r.b,e=mp(r.c);return e.$?ja(e.a):_a(kf(s(Tf,n,t,e.a)));case 
2:return _a(kf(Lf(r.a)));case 3:return _a(kf($f(r.a)));case 4:return 
_a(kf(Pf(r.a)));default:return _a(kf(f(Df,r.a,r.b)))}};function 
Dp(){return yc(N([Fd,Zf,es,f(uc,yf,va(fv)),pl,f(uc,function(){return 
_f},va(Bv)),lp,Tl(!1),va(pv),Np()]))}function Tp(){return 
yc(N([Fd,Zf,es,f(uc,yf,va(fv)),va(Iv),f(uc,function(){return 
_f},va(Bv)),lp,Tl(!0),va(pv),Np(),va(Fv)]))}function Np(){return 
f(ba,Lp,il)}function Sp(){return f(vc,{R:L,H:L},kp)}var 
Cp=Dp();Dp=function(){return Cp};var Vp=Tp();Tp=function(){return Vp};var 
Rp=Np();Np=function(){return Rp};var Up=Sp();Sp=function(){return Up};var 
Bp,Op,Pp,jp,Ip,_p,Gp,zp=e(function(r,n,t){if(1===n.$)return lt(n.a);var 
e=n.a;return 
1===t.$?lt(t.a):pt(f(r,e,t.a))}),Hp=f(De,zp(ot),pt(L)),Mp=t(function(r,n){return 
n.$?lt(n.a):r(n.a)}),Fp=t(function(r,n){r:for(;;){if(r>0){if(n.b){r-=1,n=n.b;continue 
r}return n}return n}}),Kp=e(function(r,n,t){r:for(;;){if(!t.b)return n;var 
e=t.a,u=t.b;switch(e.$){case 0:var 
a=e.a;if(a.$){i=r,o=f(r,e,n),r=i,n=o,t=u;continue r}var 
c=a.c,i=r,o=f(r,e,n);r=i,n=o,t=E(c,u);continue r;case 1:case 
2:i=r,o=f(r,e,n),r=i,n=o,t=u;continue r;case 3:var 
s=e.a;i=r,o=f(r,e,n),r=i,n=o,t=E(s,u);continue r;case 4:case 5:case 6:case 
7:default:i=r,o=f(r,e,n),r=i,n=o,t=u;continue 
r}}}),Zp=function(r){switch(r.$){case 5:return Jp(r.a);case 0:var 
n=r.a;if(n.$)return"";var e=n.c;return s(Kp,t(function(r,n){return 
E(n,Zp(r))}),"",e);case 1:return f(At,"\n",f(Te,function(r){return 
Jp(r.b)},r.a));case 2:return f(At,"\n",f(Te,Jp,r.b));case 3:return 
f(At,"\n",f(Te,Zp,e=r.a));case 4:return Jp(r.b);case 6:var u=r.b;return 
f(At,"\n",$s(N([f(Te,Jp,f(Te,function(r){return 
r.Y},r.a)),$s(f(Te,Te(Jp),u))])));case 7:return 
r.a.cV;default:return""}},Jp=function(r){return 
s(Lt,Xp,"",r)},Xp=t(function(r,n){switch(r.$){case 7:return E(n,r.a);case 
8:return n+" ";case 6:return E(n,r.a);case 1:case 2:return 
E(n,Jp(r.c));case 0:var e=r.a;if(e.$)return n;var u=e.c;return 
s(Kp,t(function(r,n){return E(n,Zp(r))}),n,u);case 4:case 3:default:return 
E(n,Jp(r.a))}}),Yp=a(function(r,n,t,e,u){var a=e;return 
f(Mp,function(e){return f(Uc,function(r){return 
r(e)},s(a,r,n,t))},Hp(u))}),Qp=function(r){return 
N([r])},Wp=e(function(r,n,t){var e=f(em,r,n);return 
e.$?t:f(ot,e.a,t)}),rm=t(function(r,n){return 
f(vb,nm(r),n)}),nm=function(r){return function(n){switch(n.$){case 4:var 
e=n.a;return ht(f(Uc,function(n){return 
r.db({cZ:n,dj:e,dI:Jp(u)})},f(um,r,u=n.b)));case 5:var u;return 
ht(f(Uc,r.dE,f(um,r,u=n.a)));case 0:var a=n.a;return 
a.$?$t:ht(l(tm,r,a.a,a.b,a.c));case 1:return 
ht(f(Uc,r.ed,Hp(f(Te,function(n){var t=n.a;return f(Uc,function(r){return 
f(Vf,t,r)},f(um,r,n.b))},c=n.a))));case 2:var c=n.b;return 
ht(f(Uc,r.dD(n.a),Hp(f(Te,um(r),c))));case 7:return ht(pt(r.c_(n.a)));case 
8:return ht(pt(r.d6));case 3:return 
ht(f(Uc,r.cU,Hp(f(rm,r,n.a))));default:var 
i=n.a,o=n.b,b=Hp(f(Te,function(n){var t=n.Y;return 
f(Uc,We(n.at),f(um,r,t))},i)),d=f(Uc,function(n){return 
r.d0(Qp(r.d2(f(Te,function(n){return 
f(r.d1,n.a,n.b)},n))))},b),v=function(n){return 
f(Uc,r.d2,f(Uc,Ct(t(function(n,t){return f(r.d$,f(hf,function(r){return 
r.at},ed(f(Fp,n,i))),t)})),Hp(f(Te,um(r),n))))},p=Hp(f(Te,v,o));return 
ht(s(zp,t(function(n,t){return r.dZ(f(ot,n,function(n){return 
ei(n)?L:N([r.d_(n)])}(t)))}),d,p))}}},tm=u(function(r,n,t,e){return 
b(Yp,n,t,e,r.dd,f(rm,r,e))}),em=t(function(r,n){switch(n.$){case 4:return 
ht(f(Uc,r.dX,f(um,r,n.a)));case 3:return ht(f(Uc,r.c4,f(um,r,n.a)));case 
5:return ht(f(Uc,r.dW,f(um,r,n.a)));case 2:var t=n.a,e=n.b;return 
ht(pt(r.de({bu:Jp(n.c),cy:t,cG:e})));case 7:return ht(pt(r.d5(n.a)));case 
6:return ht(pt(r.c$(n.a)));case 1:var u=n.a;return 
e=n.b,ht(f(Mp,function(n){return 
pt(f(r.dk,{bH:u,cG:e},n))},f(um,r,n.c)));case 8:return 
ht(pt(r.c9));default:var a=n.a;return 
a.$?$t:ht(l(tm,r,a.a,a.b,a.c))}}),um=t(function(r,n){return 
Hp(s(De,Wp(r),L,n))}),am=t(function(r,n){return 
Hp(f(rm,r,n))}),cm=e(function(r,n,t){return{$:0,a:r,b:n,c:t}}),im=e(function(r,n,t){return 
s(cm,r,n,{$:1,a:t})}),om=im("a"),fm=t(function(r,n){return{$:1,a:r,b:n}}),sm=fm,lm=function(r){return 
f(sm,"alt",r)},bm=im("blockquote"),dm={$:0},vm=e(function(r,n){return 
s(cm,r,n,dm)}),pm=vm("br"),mm=im("code"),gm=im("del"),hm=im("em"),$m=im("h1"),wm=im("h2"),ym=im("h3"),qm=im("h4"),xm=im("h5"),Am=im("h6"),km=vm("hr"),Em=function(r){return 
f(sm,"href",r)},Lm=vm("img"),Dm=im("li"),Tm=im("ol"),Nm=t(function(r,n){return 
1===r.$?n.$?lt(f(ot,r.a,n.a)):pt(n.a):pt(r.a)}),Sm=t(function(r,n){return 
ei(n)?"<"+r+">":"<"+r+" "+function(r){return f(At," 
",f(Te,function(r){return 
r.ai+'="'+r.aP+'"'},r))}(n)+">"}),Cm=im("p"),Vm=im("pre"),Rm=function(r){return 
f(sm,"src",r)},Um=im("strong"),Bm=im("table"),Om=im("tbody"),Pm=im("td"),jm=function(r){return{$:1,a:r}},Im=im("th"),_m=im("thead"),Gm=im("tr"),zm=im("ul"),Hm=function(r){return 
N(r.trim().split(/\s+/g))},Mm={cU:bm(L),c_:function(r){var 
n,t,e=r.cV,u=!(t=f(cu,Hm,r.di)).$&&t.a.b?N([(n="language-"+t.a.a,f(sm,"className",n))]):L;return 
f(Vm,L,N([f(mm,u,N([jm(e)]))]))},c$:function(r){return 
f(mm,L,N([jm(r)]))},c4:function(r){return 
f(hm,L,r)},c9:f(pm,L,L),db:function(r){var n=r.cZ;switch(r.dj){case 
0:return f($m,L,n);case 1:return f(wm,L,n);case 2:return f(ym,L,n);case 
3:return f(qm,L,n);case 4:return f(xm,L,n);default:return 
f(Am,L,n)}},dd:(Bp=L,Pp=f(Te,function(r){return 
r},Bp),Op=s(Lt,t(function(r,n){return e(function(t,e,u){return 
f(Nm,s(r,t,e,u),s(n,t,e,u))})}),e(function(){return 
lt(L)}),Pp),e(function(r,n,e){return 
f(Bc,function(e){if(e.b){if(e.b.b)return"oneOf failed parsing this 
value:\n    "+f(Sm,r,n)+"\n\nParsing failed in the following 2 
ways:\n\n\n"+f(At,"\n\n",f(Ct,t(function(r,n){return"("+xt(r+1)+") 
"+n}),e))+"\n";var u=e.a;return"Problem with the given 
value:\n\n"+f(Sm,r,n)+"\n\n"+u+"\n"}return"Ran into a oneOf with no 
possibilities!"},s(Op,r,n,e))})),de:function(r){var n=r.cG;if(n.$)return 
f(Lm,N([Rm(r.cy),lm(r.bu)]),L);var t,e=n.a;return 
f(Lm,N([Rm(r.cy),lm(r.bu),(t=e,f(sm,"title",t))]),L)},dk:t(function(r,n){return 
f(om,N([Em(r.bH)]),n)}),dD:t(function(r,n){return 
f(Tm,1===r?N([(t=r,f(sm,"start",xt(t)))]):L,f(Te,function(r){return 
f(Dm,L,r)},n));var t}),dE:Cm(L),dW:function(r){return 
f(gm,L,r)},dX:function(r){return 
f(Um,L,r)},dZ:Bm(L),d_:Om(L),d$:function(){return 
Pm(L)},d0:_m(L),d1:function(){return 
Im(L)},d2:Gm(L),d5:jm,d6:f(km,L,L),ed:function(r){return 
f(zm,L,f(Te,function(r){return f(Dm,L,r.b)},r))}},Fm=function(r){return 
r.$?$t:ht(r.a)},Km=e(function(r,n,t){return E(f(eb,r*n," 
"),t)}),Zm=t(function(r,n){if(n.b){if(n.b.b){u=n.a;var e=n.b;return 
s(Lt,t(function(n,t){return E(n,E(r,t))}),u,e)}var u;return 
n.a}return""}),Jm=function(r){return"</"+r+">"},Xm=e(function(r,n,t){return 
f(At,n,f(kt,r,t))}),Ym=f(Tc,f(Xm,"&","&amp;"),f(Tc,f(Xm,"<","&lt;"),f(Xm,">","&gt;"))),Qm=function(r){return 
r.b},Wm=f(cb,t(function(r,n){return'"'===r?n+'\\"':E(n,Sc(r))}),""),rg=f(cb,t(function(r,n){return 
Ut(r)?n+"-"+Sc(x(r.toLowerCase())):E(n,Sc(r))}),""),ng=t(function(r,n){return 
rg(r)+'="'+Wm(n)+'"'}),tg=function(r){switch(r){case"className":return"class";case"defaultValue":return"value";case"htmlFor":return"for";default:return 
r}},eg=t(function(r,n){var t=n.a,e=n.b,u=n.c;switch(r.$){case 0:return 
q(t,e,f(ot,f(ng,r.a,a=r.b),u));case 1:if("className"===r.a)return 
q(f(ot,a=r.b,t),e,u);var a=r.b;return 
q(t,e,f(ot,f(ng,tg(c=r.a),a),u));case 2:var c=r.a;return 
r.b?q(t,e,f(ot,rg(tg(c)),u)):n;case 3:return 
a=r.b,q(t,e,f(ot,f(ng,tg(c=r.a),function(r){return f(qt,0,r)}(a)),u));case 
4:return a=r.b,q(t,f(ot,Wm(r.a)+": "+Wm(a),e),u);default:return 
n}}),ug=t(function(r,n){return r.b?f(ot,f(ng,"class",f(Zm," 
",r)),n):n}),ag=t(function(r,n){return r.b?f(ot,f(ng,"style",f(Zm,"; 
",r)),n):n}),cg=t(function(r,n){return"<"+f(At," 
",f(ot,r,(t=n,e=s(Lt,eg,q(L,L,L),t),f(ag,e.b,f(ug,e.a,e.c)))))+">";var 
t,e}),ig=e(function(r,n,t){r:for(;;){if(!n.b){var e=t.N;if(e.b){var 
u=e.a,a=e.b,c=r,i=u.b,o=A(t,{p:t.p-1,r:f(ot,f(r,t.p-1,Jm(l=u.a)),t.r),N:a});r=c,n=i,t=o;continue 
r}return 
t}if(n.a.$)c=r,i=a=n.b,o=A(t,{r:f(ot,f(r,t.p,Ym(n.a.a)),t.r)}),r=c,n=i,t=o;else{var 
s=n.a,l=s.a,b=s.b,d=s.c;a=n.b;switch(d.$){case 
0:c=r,i=a,o=A(t,{r:f(ot,f(r,t.p,f(cg,l,b)),t.r)});r=c,n=i,t=o;continue 
r;case 
1:c=r,i=d.a,o=A(t,{p:t.p+1,r:f(ot,f(r,t.p,f(cg,l,b)),t.r),N:f(ot,y(l,a),t.N)}),r=c,n=i,t=o;continue 
r;default:c=r,i=f(Te,Qm,d.a),o=A(t,{p:t.p+1,r:f(ot,f(r,t.p,f(cg,l,b)),t.r),N:f(ot,y(l,a),t.N)}),r=c,n=i,t=o;continue 
r}}}}),og=t(function(r,n){var 
t=r?"\n":"",e={p:0,r:L,N:L},u=r?Km(r):Da(pe);return 
f(Zm,t,s(ig,u,N([n]),e).r)}),fg=function(r){return 
kr(s(Lt,t(function(r,n){return 
s(Lr,r.a,r.b,n)}),{},r))},sg=kr,lg=e(function(r,n,t){var 
e,u=n,a=r.aq,c=f(ju,L,f(cu,function(r){return 
N([y("format",sg("org.matrix.custom.html")),y("formatted_body",sg(r))])},(e=t,f(hf,function(r){return 
p(r,"<p>"+e+"</p>")?$t:ht(r)},f(cu,f(Tc,Te(og(0)),At("")),f(ju,$t,f(cu,f(Tc,am(Mm),Fm),s(Tc,hp,Fm,e))))))));return 
f(Fu,r,{cV:gf(fg(E(N([y("msgtype",sg("m.text")),y("body",sg(t))]),c))),dm:"PUT",dF:L,dG:N(["rooms",u,"send","m.room.message",xt(a)]),dM:de(0)})}),bg=e(function(r,n,t){var 
e=n;return f(Ne,function(){return 
s(lg,r,e,t)},f(so,r,e))}),dg=e(function(r,n,t){return 
y(s(lo(r)?lg:bg,r,n,t),A(e=r,{aq:e.aq+1}));var e}),vg=function(r){return 
r.aO},pg=t(function(r,n){return 
f(Fu,r,{cV:gf(fg(N([y("displayname",sg(n))]))),dm:"PUT",dF:L,dG:N(["profile",vg(r),"displayname"]),dM:de(0)})}),mg=function(r){var 
n=r;return 
2===n.I?A(n,{I:0}):n},gg=kr,hg=(Ip=sg,function(r){Gr[r]&&B(3)}(jp="storeSession"),Gr[jp]={e:tn,u:Ip,a:function(r){var 
n=[],t=Gr[r].u,u=_r(0);return 
Gr[r].b=u,Gr[r].c=e(function(r,e){for(;e.b;e=e.b)for(var 
a=n,c=Er(t(e.a)),i=0;a.length>i;i++)a[i](c);return 
u}),{subscribe:function(r){n.push(r)},unsubscribe:function(r){var 
t=(n=n.slice()).indexOf(r);0>t||n.splice(t,1)}}}},Kr(jp)),$g=function(r){return 
hg((e=(n=r).aC,u=n.aq,a=n.aO,c=n.U,f(qt,0,fg(N([y("homeserverUrl",sg(n.bT)),y("kind",sg((t=e,t?"user":"guest"))),y("txnId",gg(u)),y("userId",sg(a)),y("accessToken",sg(c))])))));var 
n,t,e,u,a,c},wg=t(function(r,n){return 
A(n,r.$?{ai:r.a}:{V:r.a})}),yg=function(r){return{$:5,a:r}},qg=function(r){return{$:0,a:r}},xg=function(r){return{$:1,a:r}},Ag=function(r){var 
n=r.cK,t=r.b9;return 
fg(N([y("type",sg("m.login.password")),y("identifier",fg(N([y("type",sg("m.id.user")),y("user",sg(n))]))),y("password",sg(t)),y("initial_device_display_name",sg("Cactus 
Comments"))]))},kg=function(r){var n=r.bT;return 
jo({cV:gf(Ag({b9:r.b9,cK:r.cK})),dm:"POST",dM:f(Bo,n,1),ef:s(Mu,n,N(["login"]),L)})},Eg=t(function(r,n){return 
f(He,f(ze,ou,r),n)}),Lg=function(r){return r.a},Dg=t(function(r,n){var 
t,e=r,u=e.Q;return u.$?f(Ne,function(r){return 
f(Eg,xg,kg({bT:r,b9:e.am,cK:Lg(n)}))},f(Eg,qg,(t="https://"+n.b+"/.well-known/matrix/client",Pu({cV:Oi,da:L,dm:"GET",dL:Bu(function(r){if(4===r.$){var 
n=r.b,e=f(tu,f(le,function(r){return 
f(lv,"/",r)?f(sv,1,r):r},f(Xe,"m.homeserver",f(Xe,"base_url",ru))),n);return 
e.$?lt(_t(e.a)):pt(e.a)}return lt("Failed getting 
"+t)}),d7:$t,ef:t})))):f(Eg,xg,kg({bT:u.a,b9:e.am,cK:Lg(n)}))}),Tg=t(function(r,n){var 
t,e=r,u=function(r){return q(r,vo,$t)};switch(n.$){case 1:return 
u(A(e,{am:n.a}));case 2:return u(A(e,{Q:ht(n.a)}));case 0:var a=n.a;return 
u(A(e,{aO:io(a),bs:a}));case 3:return u(A(e,{I:2}));case 4:var 
c=n.a;return q(A(e,{I:1}),f(Me,yg,f(Dg,e,c)),$t);default:if(1===n.a.$){var 
i=n.a.a;return 
u(A(e,{Q:(t=y(i,e.Q),t.a.$||1!==t.b.$?e.Q:ht("")),a5:ht(i),I:0}))}return 
q(oo,vo,ht(n.a.a))}}),Ng=t(function(r,n){if(n.$){var 
e=n.a,u=t(function(r,n){return 
0>h(Dt(af(n.C)),e.an)?f(Me,f(rf,r,0),f(sf,r,n)):vo});return 
f(df,Ie,function(){switch(r.$){case 0:return 
y(A(e,{a8:r.a}),f(ju,vo,s(bf,t(function(r,n){return 
f(Me,f(rf,r,1),f(ff,r,n))}),e.be,e.M)));case 1:var n=r.a;return 
y(A(e,{A:f(cf,function(r){return!p(r.X,n)},e.A)}),vo);case 5:return 
y(A(e,{K:f(wg,r.a,e.K)}),vo);case 6:if(r.a.$){var a=r.a.a;return 
y(A(e,{A:f(ef,e.A,a.a+" "+a.b)}),vo)}var c=r.a.a,i=c.a;return 
y(A(e,{M:ht(m=c.b),be:ht(i)}),Fe(N([$g(i),f(u,i,m)])));case 
7:if(r.c.$){var o=r.c.a;return y(A(e,{A:f(ef,e.A,o.a+" 
"+o.b)}),vo)}i=r.a;var l=r.b,b=r.c.a,d=f(cu,function(r){return 
s(mf,r,l,b)},e.M),v=A(e,{af:e.af||ei(b.bA)&&!l,M:d});return 
y(v,(D=y(d,v.af)).a.$||D.b?vo:f(u,i,D.a.a));case 8:i=r.a;var m=r.b;return 
y(A(e,{an:e.an+e.z.ba}),f(Me,f(rf,i,0),f(sf,i,m)));case 9:var 
g=s(dg,i=r.a,r.b,e.K.V),h=g.a,$=g.b;return 
y(A(e,{K:uf(e.K),be:ht($)}),Fe(N([f(Me,tf(i),function(){var r;return 
1===lo(i)?h:f(Ne,function(){return 
h},f(pg,i,""===(r=e.K).ai?"Anonymous":r.ai))}()),$g($)])));case 
10:if(r.b.$){var w=r.b.a;return y(A(e,{A:f(ef,e.A,w.a+" 
"+w.b)}),vo)}return 
y(e,f(ju,vo,f(cu,f(Tc,ff(i=r.a),Me(f(rf,i,1))),e.M)));case 4:var 
q=f(Tg,e.Z,r.a),x=q.a,k=q.b,E=f(ju,vo,f(cu,function(r){return 
f(Me,_e,f(Ne,function(){return 
f(Se,We(r),f(Fi,r,e.z.dO))},f(so,r,e.z.dO)))},$=q.c)),L=f(ju,vo,f(cu,$g,$));return 
y(A(e,{Z:x,be:$.$?e.be:$}),Fe(N([f(lf,nf,k),E,L])));case 2:return 
y(A(e,{Z:mg(e.Z)}),vo);default:return y(e,f(Me,_e,f(Ne,function(r){return 
f(Se,We(r),f(Fi,r,e.z.dO))},Io(e.z.aR))))}var D}())}return 
y(n,vo)}),Sg=function(r){return{$:5,a:r}},Cg={$:3},Vg=t(function(r,n){return{$:9,a:r,b:n}}),Rg={$:2},Ug=t(function(r,n){return{$:8,a:r,b:n}}),Bg=on("button"),Og=t(function(r,n){return 
f(bn,r,sg(n))}),Pg=Og("className"),jg=on("div"),Ig=Te(mn(Ae)),_g=function(r){return 
jg(Ig(r))},Gg=function(r){return 
r.F},zg=fn,Hg=sn,Mg=t(function(r,n){return 
f(Hg,r,{$:0,a:n})}),Fg=function(r){return 
f(Mg,"click",de(r))},Kg=an,Zg=Kg,Jg=function(r){return{$:0,a:r}},Xg=function(r){return{$:1,a:r}},Yg=on("a"),Qg=function(r){return 
Yg(Ig(r))},Wg=kr,rh=t(function(r,n){return 
f(bn,r,Wg(n))}),nh=rh("disabled"),th=function(r){return 
f(Og,"href",/^javascript:/i.test((n=r).replace(/\s/g,""))?"":n);var 
n},eh=on("input"),uh=Og("type"),ah=Og("value"),ch=t(function(r,n){return 
f(eh,E(N([uh("text"),ah(r)]),n),L)}),ih=function(r){return 
1===r.aC},oh=Og("htmlFor"),fh=ln,sh=N([f(fh,"property","clip rect(1px, 
1px, 1px, 
1px)"),f(fh,"position","absolute"),f(fh,"height","1px"),f(fh,"width","1px"),f(fh,"overflow","hidden"),f(fh,"margin","-1px"),f(fh,"padding","0"),f(fh,"border","0")]),lh=on("label"),bh=function(r){return 
lh(Ig(r))},dh=on("span"),vh=function(r){return 
dh(Ig(r))},ph=u(function(r,n,t,e){return 
f(vh,L,N([f(bh,f(ot,oh(r),E(sh,n)),N([f(zg,Ae,t)])),e]))}),mh=function(r){return 
s(Gu,"https://matrix.to",N(["#",zu(r)]),L)},gh=function(r){return 
y(r,!0)},hh=t(function(r,n){return 
f(Hg,r,{$:1,a:n})}),$h=t(function(r,n){return 
s(De,Xe,n,r)}),wh=f($h,N(["target","value"]),ru),yh=function(r){return 
f(hh,"input",f(le,gh,f(le,r,wh)))},qh=Og("placeholder"),xh=on("textarea"),Ah=e(function(r,n,t){var 
e,u=(e=y(f(cu,ih,r),f(cu,vg,r))).a.$||!e.a.a||e.b.$?"Post":"Post as 
"+e.b.a,a=p(r,$t)||!ge(t),c=E(N([Pg("cactus-button"),Pg("cactus-send-button"),nh(a)]),f(ju,L,f(cu,Qp,f(cu,Fg,n))));return 
f(Bg,c,N([Zg(u)]))}),kh=t(function(r,n){var 
t,e=r,u=n.be,a=n.dO,c=n.a4,i=n.aV,o=n.$7,b=n.dS,d=n.dl,v=s(Ah,u,n.dQ,e.V),p=f(ju,!0,f(cu,f(Tc,ih,Oa),u))?f(_g,N([Pg("cactus-editor-name")]),N([l(ph,"Name",L,Zg("Name"),f(ch,e.ai,N([qh("Name"),yh(f(Tc,Xg,o))])))])):Zg(""),m=c?function(r){var 
n=r.b0,t=r.b1,e=r.be,u=f(Bg,N([Pg("cactus-button"),Pg("cactus-logout-button"),Fg(t)]),N([Zg("Log 
out")])),a=f(Bg,N([Pg("cactus-button"),Pg("cactus-login-button"),Fg(n)]),N([Zg("Log 
in")]));return 
e.$?a:ih(e.a)?u:a}({b0:b,b1:d,be:u}):f(Qg,N([th(mh(a))]),N([f(Bg,N([Pg("cactus-button")]),N([Zg("Log 
in")]))])),g=function(r){return l(ph,"Comment Editor",L,Zg("Comment 
Editor"),f(xh,N([Pg("cactus-editor-textarea"),ah(e.V),yh(f(Tc,Jg,o)),qh("Add 
a 
comment"),nh(!r)]),L))},h=f(_g,N([Pg("cactus-editor-buttons")]),N([m,v]));return 
f(_g,N([Pg("cactus-editor")]),N((t=y(c,i)).a?t.b?[g(!0),f(_g,N([Pg("cactus-editor-below")]),N([p,h]))]:[g(f(ju,!1,f(cu,ih,u))),f(_g,N([Pg("cactus-editor-below")]),N([h]))]:t.b?[g(!0),f(_g,N([Pg("cactus-editor-below")]),N([p,h]))]:[f(Qg,N([th(mh(a)),Pg("cactus-button"),Pg("cactus-matrixdotto-only")]),N([Zg("Comment 
using a Matrix client")]))]))}),Eh=t(function(r,n){return 
f(dn,function(r){return/^(on|formAction$)/i.test(r)?"data-"+r:r}(r),vn(n))}),Lh=dn("class"),Dh=dn("d"),Th=f(ze,Eh,yt("aria-")),Nh=Th("errormessage"),Sh=on("p"),Ch=function(r){return 
Sh(Ig(r))},Vh=cn("http://www.w3.org/2000/svg"),Rh=Vh("path"),Uh=Vh("svg"),Bh=dn("viewBox"),Oh=function(r){var 
n,t=r.X,e=r.aF,u=f(Ch,N([Pg("cactus-error-text")]),N([Zg(" Error: 
"+e)])),a=f(Bg,N([Pg("cactus-error-close"),f(Eh,"aria-label","close"),Fg((n=t,{$:1,a:n}))]),N([f(Uh,N([Bh("0 
0 20 
20"),Lh("cactus-error-close-icon"),f(fh,"fill","currentColor")]),N([f(Rh,N([f(fh,"fill-rule","evenodd"),Dh("M4.293 
4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 
10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 
01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 
010-1.414z"),f(fh,"clip-rule","evenodd")]),L)]))]));return 
f(_g,N([Pg("cactus-error"),Nh(e)]),N([a,u]))},Ph=function(r){return{$:2,a:r}},jh=function(r){return{$:1,a:r}},Ih=function(r){return{$:0,a:r}},_h=function(r){return{$:4,a:r}},Gh={$:3},zh=f(Bg,N([Pg("cactus-login-close"),f(Eh,"aria-label","close"),Fg(Gh)]),N([f(Uh,N([Bh("0 
0 20 
20"),Lh("cactus-login-close-icon"),f(fh,"fill","currentColor")]),N([f(Rh,N([f(fh,"fill-rule","evenodd"),Dh("M4.293 
4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 
10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 
01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 
010-1.414z"),f(fh,"clip-rule","evenodd")]),L)]))])),Hh=f(Je,function(r){return 
f(me,"cactus-login-form-wrapper",r)?de(Gh):Zu("Ignoring click event. 
Didn't hit 
wrapper.")},f($h,N(["target","className"]),ru)),Mh=function(r){if(r.$){var 
n=r.a;return n.a+": "+n.b}return"Could not find homeserver: 
"+r.a},Fh=on("h3"),Kh=function(r){return 
Fh(Ig(r))},Zh=on("h4"),Jh=function(r){return 
Zh(Ig(r))},Xh=function(r){return 
r?"true":"false"},Yh=f(ze,Th("invalid"),Xh),Qh=Th("label"),Wh=rh("required"),r$=function(r){var 
n,t=r.ai,e=r.aP,u=r.ay,a=r.aG,c=r.au,i=r.az;return 
f(_g,N([Pg("cactus-login-field")]),N([f(bh,N([Pg("cactus-login-label")]),N([Zg(t)])),f(ch,e,E(c,N([qh(u),yh(a),Wh(!0),Qh(t),Yh((n=i,1!==n.$))]))),f(ju,Zg(""),f(cu,f(Tc,Zg,f(Tc,Qp,Ch(N([Pg("cactus-login-error")])))),i))]))},n$=t(function(r,n){var 
t,e=r,u=r$({au:L,ay:"@alice:example.com",az:(t=e.aO,1===t.$?ht(t.a):$t),aG:Ih,ai:"User 
ID",aP:e.bs}),a=f(Kh,N([Pg("cactus-login-title")]),N([Zg("Log in using 
Matrix")])),c=f(Bg,E(N([Pg("cactus-button"),Pg("cactus-login-credentials-button"),nh(!!e.I)]),f(eu,L,f(Uc,f(Tc,_h,f(Tc,Fg,Qp)),e.aO))),N([Zg(function(){switch(e.I){case 
0:return"Log in";case 1:return"Logging 
in...";default:return""}}())])),i=r$({au:N([uh("password")]),ay:"â€¢â€¢â€¢â€¢â€¢â€¢â€¢â€¢â€¢â€¢â€¢â€¢",az:$t,aG:jh,ai:"Password",aP:e.am}),o=f(ju,Zg(""),f(cu,f(Tc,Mh,f(Tc,Zg,f(Tc,Qp,Ch(N([Pg("cactus-login-error")]))))),e.a5)),s=f(ju,Zg(""),f(cu,function(r){return 
r$({au:L,ay:"https://matrix.cactus.chat:8448",az:$t,aG:Ph,ai:"Homeserver 
URL",aP:r})},e.Q)),l=f(Jh,N([Pg("cactus-login-credentials-title")]),N([Zg("Or 
type in your 
credentials")])),b=f(_g,N([Pg("cactus-login-credentials")]),N([l,u,i,s,c,o])),d=f(Jh,N([Pg("cactus-login-client-title")]),N([Zg("Use 
a Matrix 
client")])),v=f(Qg,N([th(mh(n)),Pg("cactus-button"),Pg("cactus-matrixdotto-button")]),N([Zg("Log 
in")])),p=f(_g,N([Pg("cactus-login-client")]),N([d,v]));return 
2===e.I?Zg(""):f(jg,N([Pg("cactus-login-form-wrapper"),f(Mg,"click",Hh)]),N([f(_g,N([Pg("cactus-login-form")]),N([zh,a,p,b]))]))}),t$=e(function(r,n,e){return 
f(cu,function(r){return r.ae},ed(f(Hi,f(Tc,function(r){return 
r._},f(Tc,zi,te(-1))),f(cf,function(r){return 
1>h(zi(r._),zi(n))},f(Te,Qm,f(cf,f(Tc,ae,Qt(e)),(u=r,s(Lt,t(function(r,n){return 
1===r.$?f(ot,y(r.a,r.b),n):n}),L,u))))))));var 
u}),e$=dn("datetime"),u$={$:7},a$=function(r){return 
r>11?"pm":"am"},c$=function(r){switch(r){case 0:return"January";case 
1:return"February";case 2:return"March";case 3:return"April";case 
4:return"May";case 5:return"June";case 6:return"July";case 
7:return"August";case 8:return"September";case 9:return"October";case 
10:return"November";default:return"December"}},i$=function(r){switch(f(Jl,100,r)){case 
11:case 12:case 13:return"th";default:switch(f(Jl,10,r)){case 
1:return"st";case 2:return"nd";case 
3:return"rd";default:return"th"}}},o$=function(r){switch(r){case 
0:return"Monday";case 1:return"Tuesday";case 2:return"Wednesday";case 
3:return"Thursday";case 4:return"Friday";case 
5:return"Saturday";default:return"Sunday"}},f$=d(c(function(r,n,t,e,u,a){return{d8:u,d9:n,ea:r,ab:a,eb:e,ec:t}}),c$,f(Tc,c$,ye(3)),o$,f(Tc,o$,ye(3)),a$,i$),s$=t(function(r,n){return 
Wt(r/n)}),l$=e(function(r,n,t){for(;;){if(!t.b)return n+r;var 
e=t.a,u=t.b;if(0>h(e.bm,n))return 
n+e.b;r=r,n=n,t=u}}),b$=t(function(r,n){var t=r.b;return 
s(l$,r.a,f(s$,zi(n),6e4),t)}),d$=t(function(r,n){return 
f(Jl,24,f(s$,f(b$,r,n),60))}),v$=e(function(r,n,t){return 
r.d8(f(d$,n,t))}),p$=function(r){var 
n=f(s$,r,1440)+719468,t=(0>n?n-146096:n)/146097|0,e=n-146097*t,u=(e-(e/1460|0)+(e/36524|0)-(e/146096|0))/365|0,a=e-(365*u+(u/4|0)-(u/100|0)),c=(5*a+2)/153|0,i=c+(10>c?3:-9);return{bD:a-((153*c+2)/5|0)+1,b5:i,cP:u+400*t+(i>2?0:1)}},m$=t(function(r,n){return 
p$(f(b$,r,n)).bD}),g$=N([6,0,1,2,3,4,5]),h$=t(function(r,n){switch(f(Jl,7,f(s$,f(b$,r,n),1440))){case 
0:return 3;case 1:return 4;case 2:return 5;case 3:return 6;case 4:return 
0;case 5:return 1;default:return 2}}),$$=t(function(r,n){return 
f(ju,y(0,6),ed(f(cf,function(t){return 
p(t.b,f(h$,r,n))},f(Ct,t(function(r,n){return 
y(r,n)}),g$)))).a}),w$=t(function(r,n){switch(n){case 0:return 31;case 
1:return 
function(r){return!(f(Jl,4,r)||!f(Jl,100,r)&&f(Jl,400,r))}(r)?29:28;case 
2:return 31;case 3:return 30;case 4:return 31;case 5:return 30;case 6:case 
7:return 31;case 8:return 30;case 9:return 31;case 10:return 
30;default:return 
31}}),y$=N([0,1,2,3,4,5,6,7,8,9,10,11]),q$=t(function(r,n){switch(p$(f(b$,r,n)).b5){case 
1:return 0;case 2:return 1;case 3:return 2;case 4:return 3;case 5:return 
4;case 6:return 5;case 7:return 6;case 8:return 7;case 9:return 8;case 
10:return 9;case 11:return 10;default:return 
11}}),x$=t(function(r,n){return f(ju,y(0,0),ed(f(cf,function(t){return 
p(t.b,f(q$,r,n))},f(Ct,t(function(r,n){return 
y(r,n)}),y$))))}),A$=t(function(r,n){return 
1+f(x$,r,n).a}),k$=t(function(r,n){return 
p$(f(b$,r,n)).cP}),E$=t(function(r,n){var t,e=f(Yv,f(A$,r,n)-1,y$);return 
t=f(Te,w$(f(k$,r,n)),e),s(Lt,gt,0,t)+f(m$,r,n)}),L$=t(function(r,n){return 
f(A$,r,n)/4|0}),D$=t(function(r,n){return 
1>r?"":s(he,-r,ge(n),n)}),T$=t(function(r,n){var t=xt(n),e=r-ge(t);return 
E(f(At,"",f(Te,function(){return"0"},f(St,1,e))),t)}),N$=t(function(r,n){return 
f(Jl,1e3,zi(n))}),S$=t(function(r,n){return 
f(Jl,60,f(b$,r,n))}),C$=function(r){return 
r?r>12?r-12:r:12},V$=t(function(r,n){return 
f(Jl,60,f(s$,zi(n),1e3))}),R$=z(315576e5),U$=t(function(r,n){return 
Ci(R$*f(k$,r,n))}),B$=t(function(r,n){var t=f($$,r,f(U$,r,n));return 
1+((f(E$,r,n)+t)/7|0)}),O$=t(function(r,n){return 
xt(f(k$,r,n))}),P$=u(function(r,n,t,e){switch(e.$){case 0:return 
xt(f(A$,n,t));case 1:return E(xt(u=f(A$,n,t)),r.ab(u));case 2:return 
f(T$,2,f(A$,n,t));case 3:return r.d9(f(q$,n,t));case 4:return 
r.ea(f(q$,n,t));case 17:return xt(1+f(L$,n,t));case 18:return 
function(n){return E(xt(n),r.ab(n))}(1+f(L$,n,t));case 5:return 
xt(f(m$,n,t));case 6:return function(n){return 
E(xt(n),r.ab(n))}(f(m$,n,t));case 7:return f(T$,2,f(m$,n,t));case 8:return 
xt(f(E$,n,t));case 9:return function(n){return 
E(xt(n),r.ab(n))}(f(E$,n,t));case 10:return f(T$,3,f(E$,n,t));case 
11:return xt(f($$,n,t));case 12:return function(n){return 
E(xt(n),r.ab(n))}(f($$,n,t));case 13:return r.eb(f(h$,n,t));case 14:return 
r.ec(f(h$,n,t));case 19:return xt(f(B$,n,t));case 20:return 
function(n){return E(xt(n),r.ab(n))}(f(B$,n,t));case 21:return 
f(T$,2,f(B$,n,t));case 15:return f(D$,2,f(O$,n,t));case 16:return 
f(O$,n,t);case 22:return s(v$,r,n,t).toUpperCase();case 23:return 
za(s(v$,r,n,t));case 24:return xt(f(d$,n,t));case 25:return 
f(T$,2,f(d$,n,t));case 26:return xt(C$(f(d$,n,t)));case 27:return 
f(T$,2,C$(f(d$,n,t)));case 28:return xt(1+f(d$,n,t));case 29:return 
f(T$,2,1+f(d$,n,t));case 30:return xt(f(S$,n,t));case 31:return 
f(T$,2,f(S$,n,t));case 32:return xt(f(V$,n,t));case 33:return 
f(T$,2,f(V$,n,t));case 34:return xt(f(N$,n,t));case 35:return 
f(T$,3,f(N$,n,t));default:return e.a}var u}),j$=u(function(r,n,t,e){return 
f(At,"",f(Te,s(P$,r,t,e),n))})(f$),I$={$:25},_$={$:31},G$={$:2},z$={$:33},H$=function(r){return{$:36,a:r}},M$=f(po,0,L),F$={$:16},K$={$:13},Z$={$:3},J$=on("time"),X$=function(r){return 
J$(Ig(r))},Y$=t(function(r,n){return function(r){return 
Ro(.001*r)}(zi(n)-zi(r))}),Q$=function(r){return 
Qo(r)/86400},W$=function(r){return Qo(r)/3600},rw=function(r){return 
Qo(r)/31557600},nw=function(r){return Qo(r)/60},tw=function(r){return 
Qo(r)/604800},ew=t(function(r,n){var 
t=N([y("year",rw),y("month",f(Tc,rw,te(12))),y("week",tw),y("day",Q$),y("hour",W$),y("minute",nw),y("second",Qo)]),e=f(Y$,n,r),u=function(r){r:for(;;){if(r.b){var 
n=r.a,t=n.a,u=n.b,a=r.b;if(1>u(e)){r=a;continue r}return y(t,u)}return 
y("seconds",Qo)}}(t),a=u.a,c=Wt((0,u.b)(e)),i=1===c?a:a+"s";return 
1>Qo(e)?"just now":xt(c)+" "+i+" 
ago"}),uw=Og("title"),aw=Og("alt"),cw=on("img"),iw=t(function(r,n){return 
f(cw,f(ot,aw(r),Ig(n)),L)}),ow=function(r){return 
f(Og,"src",vn(r))},fw=t(function(r,n){return 
f(Oo,zu(r),xt(n))}),sw=Hu(N(["_matrix","media","r0"])),lw=function(r){return 
s(Tc,jt,ed,f(kt,"/",r))},bw=function(r){var n=function(r){return 
Pt(r)||f(oc,r,N([".","-",":"]))},t=f(Wa,f(Sa,Ga(pe),tc("mxc://")),ao({a2:n,bc:Ji,bm:n}));return 
Fm(f(qi,t,r))},dw=t(function(r,n){var e=bw(n),u=lw(n);return 
s(bf,t(function(n,t){return 
s(sw,r,N(["thumbnail",n,t]),N([f(fw,"width",64),f(fw,"height",64),f(Po,"method","crop")]))}),e,u)}),vw=t(function(r,n){var 
t,e=(t=f(cu,dw(r),f(hf,function(r){return r.bw},n))).$?$t:t.a;return 
f(_g,N([Pg("cactus-comment-avatar")]),N(e.$?[f(_g,N([Pg("cactus-comment-avatar-placeholder")]),L)]:[f(iw,"user 
avatar image",N([ow(e.a)]))]))}),pw=on("audio"),mw=function(r){return 
pw(Ig(r))},gw=rh("controls"),hw=t(function(r,n){var e=bw(n),u=lw(n);return 
s(bf,t(function(n,t){return 
s(sw,r,N(["download",n,t]),L)}),e,u)}),$w=on("i"),ww=function(r){return 
$w(Ig(r))},yw=t(function(r,n){var t=f(hw,r,n.ef);if(t.$)return 
f(Ch,L,N([f(ww,L,N([Zg("Error: Could not parse audio file URL")]))]));var 
e=t.a;return 
f(mw,N([Pg("cactus-message-audio"),ow(e),gw(!0)]),L)}),qw=dn("rel"),xw=t(function(r,n){var 
t=f(hw,r,n.ef);if(t.$)return f(Ch,L,N([f(ww,L,N([Zg("Error: Could not 
parse file URL")]))]));var e=t.a;return 
f(Qg,N([Pg("cactus-message-file"),qw("noopener"),th(e)]),N([Zg("ðŸ“Ž 
Download "+n.cV)]))}),Aw=f(da,function(r){return 7===ge(r)?Ga(r):Ia("Hex 
color code should have 7 
characters.")},Ua(f(Sa,f(Sa,Yf(f(Xa,"#",Ds("#"))),La(Ic)),Qi))),kw=t(function(r,n){return 
s(Lt,t(function(n,t){switch(n.a){case"width":case"height":case"alt":case"title":return 
f(ot,n,t);case"src":return f(ju,t,f(cu,function(r){return 
f(ot,y("src",r),t)},f(hw,r,n.b)));default:return 
t}}),L,n)}),Ew=e(function(r,n,e){switch(n){case"font":case"span":return 
function(r){return 
s(Lt,t(function(r,n){switch(r.a){case"data-mx-color":var 
t=f(qi,Aw,r.b);return 
t.$?n:f(ot,y("style","color:"+t.a),n);case"data-mx-bg-color":var 
e=f(qi,Aw,r.b);return 
e.$?n:f(ot,y("style","background:"+e.a),n);default:return 
n}}),L,r)}(e);case"a":return f(ot,y("rel","noopener"),function(r){return 
s(Lt,t(function(r,n){switch(r.a){case"name":return 
f(ot,y("name",r.b),n);case"target":return 
f(ot,y("target",r.b),n);case"href":var 
t=r.b,e=N(["https","http","ftp","mailto","magnet"]);return 
f(ju,!1,f(cu,function(r){return 
f(oc,r,e)},ed(f(kt,":",t))))?f(ot,y("href",t),n):n;default:return 
n}}),L,r)}(e));case"img":return f(kw,r,e);case"ol":return 
f(cf,f(Tc,ae,Qt("start")),e);case"code":return 
f(cf,function(r){return"class"===r.a&&"language-"===f(ye,9,r.b)},e);default:return 
L}}),Lw=t(function(r,n){return 
s(Lu,r,0,n)}),Dw=(_p=N(["font","del","h1","h2","h3","h4","h5","h6","blockquote","p","a","ul","ol","sup","sub","li","b","i","u","strong","em","strike","code","hr","br","div","table","thead","tbody","tr","th","td","caption","pre","span","img"]),s(Lt,Lw,Ji,_p)),Tw=t(function(r,n){var 
e=t(function(n,t){if(n>100)return si("");switch(t.$){case 0:return 
si(t.a);case 2:return Ma(t.a);default:var 
u=t.a,a=t.b,c=t.c;return(f(eo,u,Dw)?f(fa,u,s(Ew,r,u,a)):f(fa,"div",L))(f(Te,e(n+1),c))}});return 
f(e,0,n)}),Nw=function(r){return 
on(function(r){return"script"==r?"p":r}(r))},Sw=function(r){return 
f(Eh,r.a,r.b)},Cw=function(r){return 
f(Te,Vw,r)},Vw=function(r){switch(r.$){case 1:var n=r.c;return 
s(Nw,r.a,f(Te,Sw,r.b),Cw(n));case 0:return Kg(r.a);default:return 
Kg("")}},Rw=t(function(r,n){if(n.$){var t=n.a;return 
Cw(f(Te,Tw(r),t))}return N([f(Ch,L,N([Zg(n.a)]))])}),Uw=function(r){return 
f(dn,"height",xt(r))},Bw=function(r){return 
f(dn,"width",xt(r))},Ow=t(function(r,n){var t=f(hf,function(n){var 
t=y(f(ju,$t,f(cu,hw(r),n.cF)),n.cE);return 
t.a.$||t.b.$?$t:ht(y(t.a.a,t.b.a))},n.a1),e=y(f(hw,r,n.ef),n.a1),u=y(e,t);if(u.a.a.$)return 
f(Ch,L,N([f(ww,L,N([Zg("Error: Could not render 
image")]))]));if(u.b.$){if(u.a.b.$)return 
f(Qg,N([th(i=u.a.a.a)]),N([f(iw,n.cV,N([Pg("cactus-message-image"),ow(i)]))]));var 
a=u.a,c=a.b.a;return 
f(Qg,N([th(i=a.a.a)]),N([f(iw,n.cV,N([Pg("cactus-message-image"),ow(i),Bw(c.aQ),Uw(c.aB)]))]))}var 
i,o=u.b.a,s=o.a,l=o.b;return 
f(Qg,N([th(i=u.a.a.a)]),N([f(iw,n.cV,N([Pg("cactus-message-image"),ow(s),Bw(l.aQ),Uw(l.aB)]))]))}),Pw=on("video"),jw=function(r){return 
Pw(Ig(r))},Iw=t(function(r,n){var t=f(hw,r,n.ef);if(t.$)return 
f(Ch,L,N([f(ww,L,N([Zg('Error: The URL for video "'+n.cV+'" could not be 
decoded.')]))]));var e=t.a;return 
f(jw,N([aw(n.cV),ow(e),gw(!0),Pg("cactus-message-video")]),L)}),_w=e(function(r,n,t){switch(t.$){case 
0:var e=t.a;return f(_g,N([Pg("cactus-message-text")]),f(Rw,r,e));case 
1:if(t.a.$)return e=t.a,f(_g,N([Pg("cactus-message-text")]),f(Rw,r,e));var 
u=t.a.a;return f(_g,N([Pg("cactus-message-emote")]),N([f(Ch,L,N([Zg(n+" 
"+u)]))]));case 2:return 
e=t.a,f(_g,N([Pg("cactus-message-text")]),f(Rw,r,e));case 3:return 
f(Ow,r,t.a);case 4:return f(xw,r,t.a);case 6:return 
f(yw,r,t.a);default:return f(Iw,r,t.a)}}),Gw=c(function(r,n,t,e,u,a){var 
c,i=(c=t,f(j$(N([F$,H$("-"),G$,H$("-"),u$,H$("T"),I$,H$(":"),_$,H$(":"),z$,H$("+00:00")])),M$,c)),o=function(r){return 
f(j$(N([K$,H$(" "),Z$,H$(" "),u$,H$(" "),I$,H$(":"),_$,H$(":"),z$,H$(" 
"),F$,H$(" 
UTC")])),M$,r)}(t),l=f(ew,n,t),b="https://matrix.to/#/"+e,d=f(ju,e,f(cu,function(r){return 
f(ju,e,r.bI)},u)),v=s(_w,r,d,a);return 
f(_g,N([Pg("cactus-comment")]),N([f(vw,r,u),f(_g,N([Pg("cactus-comment-content")]),N([f(_g,N([Pg("cactus-comment-header")]),N([f(Qg,N([Pg("cactus-comment-displayname"),th(b)]),N([Zg(d)])),f(X$,N([Pg("cactus-comment-time"),uw(o),e$(i)]),N([Zg(l)]))])),f(_g,N([Pg("cactus-comment-body")]),N([v]))]))]))}),zw=u(function(r,n,t,e){var 
u=n;return f(_g,N([Pg("cactus-comments-list")]),f(Te,function(n){return 
d(Gw,r,e,n._,n.aL,f(ju,f(xu,n.aL,u.b2),f(cu,ht,s(t$,u.C,n._,n.aL))),n.ae)},f(Yv,t,af(u.C))))});Gp={Main:{init:je({dg:function(r){var 
n=function(r){return f(Uc,Uo,f(Lo,Co,r))}(r);if(n.$)return 
y({$:0,a:_t(n.a)},vo);var t=n.a,e=t.a,u=t.b;return 
y(Ie({z:e,K:Ki,A:L,af:!1,Z:oo,a8:Ci(0),M:$t,be:u,an:e.ba}),Fe(N([f(Pe,Ge,mo),f(Me,_e,function(){if(1===u.$)return 
f(Ne,function(r){return f(Se,We(r),f(Fi,r,e.dO))},Io(e.aR));var 
r=u.a;return f(Ne,function(){return 
f(Se,We(r),f(Fi,r,e.dO))},f(bo,r,e.dO))}())])))},dY:function(r){if(1===r.$){var 
n=1e3*Qo(r.a.z.br);return n>0?f(Yo,n,Ge):Wo}return 
Wo},ee:Ng,eg:function(r){if(r.$){var 
n=r.a,t=f(zg,nf,f(n$,n.Z,n.z.dO)),e=Dt(n.A)>0?f(_g,N([Pg("cactus-errors")]),f(Te,Oh,n.A)):Zg(""),u=f(kh,n.K,{aV:n.z.aV,a4:n.z.a4,dl:Cg,$7:Sg,dO:n.z.dO,dQ:s(bf,Vg,n.be,f(cu,Gg,n.M)),be:n.be,dS:Rg});return 
f(_g,N([Pg("cactus-container")]),N([e,t,u,function(){var 
r=y(n.M,n.be);if(r.a.$||r.b.$)return 
f(_g,N([Pg("spinner")]),N([f(_g,N([Pg("rect1")]),L),f(_g,N([Pg("rect2")]),L),f(_g,N([Pg("rect3")]),L),f(_g,N([Pg("rect4")]),L)]));var 
t,e=r.a.a,u=r.b.a;return 
f(_g,N([Pg("cactus-comments-container")]),N([l(zw,(t=u,t.bT),e,n.an,n.a8),n.af?Zg(""):f(_g,N([Pg("cactus-view-more")]),N([f(Bg,N([Pg("cactus-button"),Fg(f(Ug,u,e))]),N([Zg("View 
more")]))]))]))}()]))}return Oh({X:0,aF:"Bad configuration: 
"+r.a})}})(Do)(0)}},r.Elm?function r(n,t){for(var e in t)e in 
n?"init"==e?B(6):r(n[e],t[e]):n[e]=t[e]}(r.Elm,Gp):r.Elm=Gp}(this);},{}],"Yscq":[function(require,module,exports){"use 
strict";var e,t,n=require("./Main.elm");function 
o(e){e.storedSession=JSON.parse(localStorage.getItem("cactus-session"));var 
t=e.node;delete e.node,"string"==typeof 
t&&(t=document.querySelector(t)),n.Elm.Main.init({node:t,flags:e}).ports.storeSession.subscribe(function(e){return 
localStorage.setItem("cactus-session",e)})}(null===(e=document.currentScript)||void 
0===e?void 0:null===(t=e.dataset)||void 0===t?void 
0:t.defaultHomeserverUrl)&&o(Object.assign({node:document.currentScript},document.currentScript.dataset)),window.initComments=o;},{"./Main.elm":"asWa"}]},{},["Yscq"],null)}

