// Return v if v is not null/undefined/false, else return altv.
export function ifnull(v, altv) {
    if (v) return v;
    return altv;
}
function lskey(ns, k) {
    return `${ns}_${k}`;
}
// Return localstorage item of key: "<ns>_<k>"
export function getls(ns, k, vdefault) {
    return ifnull(localStorage.getItem(lskey(ns, k)), vdefault);
}
// Set localstorage item to key: "<ns>_<k>"
export function setls(ns, k, v) {
    localStorage.setItem(lskey(ns, k), v);
}

export function getlsInt(ns, k, vdefault) {
    let v = localStorage.getItem(lskey(ns, k));
    if (v == null) {
        return vdefault;
    }
    return parseInt(v, 10);
}
export function setlsInt(ns, k, v) {
    localStorage.setItem(lskey(ns, k), v.toString());
}

export function readCookie(name) {
    let cookies = document.cookie.split(";");
    for (let i=0; i < cookies.length; i++) {
        let cookie = cookies[i].trim();
        let [k,v] = cookie.split("=");
        if (k == name) {
            if (v == undefined) {
                v = "";
            }
            return v;
        }
    }
    return "";
}

export function currentSession() {
    let suserid = readCookie("userid");
    if (suserid == "") {
        return {userid: 0, username: "", sig: ""};
    }
    let username = readCookie("username");
    let sig = readCookie("sig");

    let userid = parseInt(suserid, 10);
    return {
        userid: userid,
        username: username,
        sig: sig,
    };
}

export function initPopupHandlers() {
    function onglobalclick(e) {
        // Send signal to close any open pop-up menus.
        let mm = document.querySelectorAll(".popupmenu");
        for (let i=0; i < mm.length; i++) {
            let e = new Event("globalclick");
            mm[i].dispatchEvent(e);
        }
    }
    document.addEventListener("click", onglobalclick, false);
}

// All purpose GET request function
// sreq contains the uri with query params,
// Ex. "/api/entry?id=123" or "/api/entries?userid=123"
// fmt contains "json" or undefined/null to return json object
// any other value to fmt will return plaintext object representation
// Returns [null, err] if an error occured.
// Returns [item, null] if successful, where item contains value returned from request.
export async function find(sreq, fmt) {
    try {
        let res = await fetch(sreq, {method: "GET"});
        if (!res.ok) {
            if (res.status == 404) {
                return [null, null];
            }
            let s = await res.text();
            let err = new Error(s);
            err.status = res.status;
            return [null, err];
        }
        let v;
        if (!fmt || fmt == "json") {
            v = await res.json();
        } else {
            v = await res.text();
        }
        return [v, null];
    } catch(err) {
        return [null, err];
    }
}

// All purpose POST/PUT request function
// sreq contains the uri
// method contains the http method ("POST" or "PUT")
// item contains the object to be submitted
// Returns [null, err] if an error occured.
// Returns [item, null] if successful, where item contains final object saved.
export async function submit(sreq, method, item) {
    try {
        let res = await fetch(sreq, {
            method: method,
            headers: {"Content-Type": "application/json"},
            body: JSON.stringify(item),
        });
        if (!res.ok) {
            let s = await res.text();
            let err = new Error(s);
            err.status = res.status;
            return [null, err];
        }
        let saveditem = await res.json();
        return [saveditem, null];
    } catch(err) {
        return [null, err];
    }
}

// Similar to submit(), but ignore any response object and always use POST.
// sreq contains the uri
// Returns err if an error occured, or null if successful.
export async function exec(sreq, item) {
    try {
        let res;
        if (item) {
            res = await fetch(sreq, {
                method: "POST",
                headers: {"Content-Type": "application/json"},
                body: JSON.stringify(item),
            });
        } else {
            res = await fetch(sreq, {
                method: "POST",
            });
        }
        if (!res.ok) {
            let s = await res.text();
            let err = new Error(s);
            err.status = res.status;
            return err;
        }
        return null;
    } catch(err) {
        return err;
    }
}

// All purpose DELETE request function
// sreq contains the uri
// Returns err if an error occured, or null if successful.
export async function del(sreq) {
    try {
        let res = await fetch(sreq, {method: "DELETE"});
        if (!res.ok) {
            let s = await res.text();
            let err = new Error(s);
            err.status = res.status;
            return err;
        }
        return null;
    } catch(err) {
        return err;
    }
}

export async function subscribe(sreq, fmt, respCB) {
    let [msg, err] = await find(sreq, fmt);
    if (respCB != null) {
        respCB(msg, err)
    }

    let waitms = 0;
    // If error, wait 5 secs before trying to reconnect.
    if (err != null) {
        waitms = 5000;
    }

    setTimeout(async function() {
        await subscribe(sreq, fmt, respCB);
    }, waitms);
}

export function formatnum(snum) {
    let rx = /^\s*(\d*)\.?(\d*)\s*$/;
    let matches = snum.match(rx);

    // Ex. "123.45"
    // matches: ["123.45", "123", "45"]

    if (matches == null || matches.length == 1) {
        return "0.00";
    }
    let sleft = matches[1];
    let sright = matches[2];
    if (sleft == "") {
        sleft = "0";
    }
    if (sright == "") {
        sright = "00";
    } else if (sright.length == 1) {
        sright = `${sright}0`;
    }
    return `${sleft}.${sright}`;
}

export function textToHtml(s) {
    let ss = escapehtml(s.trim()).split("\n");
    return "<p>" + ss.join("</p><p>") + "</p>";
}
// Thanks to https://stackoverflow.com/questions/2794137/sanitizing-user-input-before-adding-it-to-the-dom-in-javascript
export function escapehtml(s) {
    let m = {
        "&": "&amp;",
        "<": "&lt;",
        ">": "&gt;",
        '"': "&quot;",
        "'": "&#x27;",
        "/": "&#x2F;",
    };
    let rx = /[&<>"'/]/ig;
    return s.replace(rx, match => m[match]);
}

