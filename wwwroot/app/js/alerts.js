var alertContainer = function (id) {
    this.container = document.getElementById(id);
    this._count = {};
    this.add = function (type, msg) {
        var d = document.createElement('div');
        d.className = 'alert alert-' + type;
        var b = document.createElement('button');
        b.setAttribute('type', 'button');
        b.className = 'close';
        b.addEventListener('click', function (evt) {
            var a = this.parentElement;
            a.parentNode.removeChild(a);
        });
        var s = document.createElement('span');
        s.innerHTML = '&times;';
        b.appendChild(s);
        d.appendChild(b);
        d.appendChild(document.createTextNode(msg));
        this.container.appendChild(d);
        var c = this._count[type] == undefined ? 0 : this._count[type];
        this._count[type] = ++c;
    };
    this.clear = function () {
        this._count = {};
        while (this.container.hasChildNodes()) {
            this.container.removeChild(this.container.lastChild);
        }
    };
    this.count = function (type) {
        return this._count[type] == undefined ? 0 : this._count[type];
    }
};