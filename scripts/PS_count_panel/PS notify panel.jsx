#target photoshop
// Sends the current open documents (count + active + names) to the
// Open Files Counter panel over a local TCP socket. Works even during a
// Batch, because the panel receives it in a separate process (not through
// PS's blocked main thread). Payload format: COUNT|ACTIVE|name1|name2|...
try {
    var c = app.documents.length;
    var a = c > 0 ? app.activeDocument.name : "";
    var names = [];
    for (var i = 0; i < c; i++) names.push(app.documents[i].name);
    var payload = c + "|" + a + "|" + names.join("|");

    var s = new Socket();
    if (s.open("127.0.0.1:45678", "UTF-8")) {
        s.write(payload);
        s.close();
    }
} catch (e) {}
