# language: de
# Spec: OPTIONS-001 — Serveroptionen ohne Build-Bindung
# Jedes Szenario beschreibt das *Verhalten* einer Leitplanke und ist mit ihrer
# PG-ID getaggt. Die IDs hier müssen deckungsgleich mit `guardrails` in
# OPTIONS-001.spec.md sein (Abgleich siehe specs/README.md).
#
# Diese Datei wird nicht ausgeführt — pz-admin hat keine Testsuite.
# Durchgesetzt wird über die Kommandos und Handgriffe im Gate-Contract.

Funktionalität: Der Server sagt, welche Optionen es gibt
  Die Anwendung liest die Optionsliste vollständig aus `showoptions` und schreibt
  ausschließlich zurück, was von dort kam. Was sie nicht kennt, zeigt sie roh —
  statt es stillschweigend fallen zu lassen.

  Grundlage:
    Angenommen eine bestehende RCON-Verbindung zu einem Server
    Und die Anwendung hat kein Vorwissen über die Optionsliste dieses Servers

  @PG-01
  Szenario: Eine unbekannte Option verschwindet nicht
    Angenommen ein Server, der die Option "War" meldet
    Und die Anwendung kennt diese Option nicht
    Wenn der Optionen-Tab geöffnet wird
    Dann erscheint "War" mit ihrem aktuellen Wert
    Und sie lässt sich ändern und speichern

  @PG-01
  Szenario: Ein neuer Server-Build braucht kein neues Release
    Angenommen ein Server-Update, das zwölf neue Optionen mitbringt
    Wenn die Anwendung sich erneut verbindet
    Dann sind alle zwölf sichtbar
    Und keine davon musste im Code eingetragen werden

  @PG-02
  Szenario: Die Darstellungstabelle entscheidet nicht über Existenz
    Angenommen eine Option ohne Eintrag in der Darstellungstabelle
    Wenn die Oberfläche sie anzeigt
    Dann steht sie unter "Weitere Optionen" mit ihrem Rohnamen
    Und ihr Fehlen in der Tabelle ist ein Hinweis, keine Sperre

  @PG-03
  Szenario: Es wird nichts geschickt, was der Server nicht kennt
    Angenommen ein Server, der die Option "KickFastPlayers" nicht mehr führt
    Wenn Optionen gespeichert werden
    Dann wird für diese Option kein "changeoption" gesendet
    Und die Anwendung meldet dafür auch keinen Fehlschlag

  @PG-03
  Szenario: Nur Geändertes wird geschrieben
    Angenommen ein Server mit einhundertvierundvierzig Optionen
    Wenn genau zwei davon geändert und gespeichert werden
    Dann werden genau zwei "changeoption" gesendet
    Und danach genau ein "reloadoptions"

  @PG-04
  Szenario: Anti-Cheat ist eine Auswahl aus vier Werten
    Angenommen ein Server auf Build 42
    Wenn die Option "AntiCheatSafety" angezeigt wird
    Dann bietet sie Ban, Kick, Log und Disable an
    Und keinen Ein-Aus-Schalter

  @PG-04
  Szenario: Die alten Anti-Cheat-Schalter sind verschwunden
    Angenommen die Oberfläche auf einem Build-42-Server
    Wenn der Anti-Cheat-Bereich geöffnet wird
    Dann taucht keine Option namens "AntiCheatProtectionType1" auf
    Und auch kein Schwellenwert-Multiplikator dazu

  @PG-05
  Szenario: Ein Wert unbekannten Typs überlebt das Speichern
    Angenommen die Option "ClientCommandFilter" mit Semikolons und Sternchen im Wert
    Wenn sie unverändert gespeichert wird
    Dann liefert der Server danach denselben Wert zurück
    Und nichts wurde abgeschnitten oder umgedeutet

  @PG-05
  Szenario: Eine leere Zeichenkette ist ein Wert
    Angenommen eine Option, deren Wert leer ist
    Wenn die Oberfläche sie anzeigt
    Dann steht dort ein leeres Textfeld
    Und beim Speichern ohne Eingabe wird sie nicht gesendet

  @PG-06
  Szenario: Eine abgelehnte Option wird namentlich gemeldet
    Angenommen eine Option, die der Server beim Setzen zurückweist
    Wenn gespeichert wird
    Dann nennt die Meldung den Namen dieser Option
    Und nicht nur die Anzahl der Fehlschläge

  @PG-06
  Szenario: Teilerfolg wird als Teilerfolg gemeldet
    Angenommen fünf geänderte Optionen, von denen eine zurückgewiesen wird
    Wenn gespeichert wird
    Dann meldet die Anwendung vier erfolgreiche und eine fehlgeschlagene
    Und die vier bleiben gesetzt
