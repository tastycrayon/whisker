import{w}from"./index.418f7e58.js";import{ao as p}from"./index.7bf39313.js";const d={};function S(e){return e==="local"?localStorage:sessionStorage}function i(e,r,s){const o=s?.serializer??JSON,m=s?.storage??"local";function f(n,c){S(m).setItem(n,o.stringify(c))}if(!d[e]){const n=w(r,t=>{const a=S(m).getItem(e);a&&t(o.parse(a));{const g=l=>{l.key===e&&t(l.newValue?o.parse(l.newValue):null)};return window.addEventListener("storage",g),()=>window.removeEventListener("storage",g)}}),{subscribe:c,set:u}=n;d[e]={set(t){f(e,t),u(t)},update(t){const a=t(p(n));f(e,a),u(a)},subscribe:c}}return d[e]}const P=i("modeOsPrefers",!1),h=i("modeUserPrefers",void 0),M=i("modeCurrent",!1);function _(){const e=window.matchMedia("(prefers-color-scheme: light)").matches;return P.set(e),e}function C(e){h.set(e)}function v(e){const r=document.documentElement.classList,s="dark";e===!0?r.remove(s):r.add(s),M.set(e)}function E(){const e=document.documentElement.classList,r=localStorage.getItem("modeUserPrefers")==="false",s=!("modeUserPrefers"in localStorage),o=window.matchMedia("(prefers-color-scheme: dark)").matches;r||s&&o?e.add("dark"):e.remove("dark")}export{v as a,C as b,_ as g,i as l,M as m,E as s};