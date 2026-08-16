# language: de
# Spec: CMD-001 — Die Befehlsfläche von Build 42
# Jedes Szenario beschreibt das *Verhalten* einer Leitplanke und ist mit ihrer
# PG-ID getaggt. Die IDs hier müssen deckungsgleich mit `guardrails` in
# CMD-001.spec.md sein (Abgleich siehe specs/README.md).
#
# Diese Datei wird nicht ausgeführt — pz-admin hat keine Testsuite.

Funktionalität: Was der Server kann, kann auch die Oberfläche
  Die neunzehn fehlenden Admin-Befehle bekommen einen Einstiegspunkt, die drei
  Karteileichen verschwinden. Wer etwas auslöst, das sich nicht zurücknehmen lässt,
  bekommt vorher gesagt, was passiert.

  Grundlage:
    Angenommen eine bestehende RCON-Verbindung zu einem Build-42-Server
    Und ein angemeldeter Nutzer mit Adminrechten

  @PG-07
  Szenario: Ein Befehl, den es nicht mehr gibt, wird nicht vorgeschlagen
    Angenommen das Terminal mit seiner Autovervollständigung
    Wenn "gr" getippt wird
    Dann wird "grantadmin" nicht vorgeschlagen
    Und auch "removeadmin" und "replay" stehen nicht in der Liste

  @PG-07
  Szenario: Die Liste nennt ihren Stand
    Angenommen die Befehlsliste im Quelltext
    Wenn sie gelesen wird
    Dann steht dabei, gegen welchen Build sie gepflegt ist
    Und ein späterer Abgleich mit dem Wiki ist ohne Raten möglich

  @PG-08
  Szenario: Die Oberfläche baut keine Kommandozeile
    Angenommen ein Dialog, der einen Spieler aus einem Safehouse wirft
    Wenn er abgeschickt wird
    Dann ruft er eine benannte Funktion mit Spielernamen auf
    Und setzt selbst keine Zeichenkette mit Anführungszeichen zusammen

  @PG-08
  Szenario: Rohzugriff gibt es genau an einer Stelle
    Angenommen die gesamte Oberfläche
    Wenn nach Aufrufen des rohen Befehlsversands gesucht wird
    Dann findet sich genau einer, im Terminal
    Und dort ist er ausdrücklich gewollt

  @PG-09
  Szenario: Eine freundliche Antwort ist kein Erfolg
    Angenommen ein Befehl auf einen Spieler, den es nicht gibt
    Wenn der Server mit einem Hinweistext und ohne Fehler antwortet
    Dann meldet die Anwendung einen Fehlschlag
    Und keinen Erfolg

  @PG-09
  Szenario: Ein Befehl über mehrere Spieler kann teilweise gelingen
    Angenommen drei ausgewählte Spieler, von denen einer offline geht
    Wenn der Befehl ausgeführt wird
    Dann meldet die Anwendung zwei Erfolge und einen Fehlschlag
    Und nennt den betroffenen Namen

  @PG-09
  Szenario: Fernwirkung nimmt die lange Befehlsform
    Angenommen der Wunsch, einen anderen Spieler unverwundbar zu machen
    Wenn die Oberfläche den Befehl absetzt
    Dann benutzt sie die Form für andere Spieler
    Und nicht die Kurzform, die auf die Konsole selbst zielt

  @PG-10
  Szenario: Zombies entfernen fragt zurück und sagt was
    Angenommen der Knopf zum Entfernen der Zombies
    Wenn er gedrückt wird
    Dann nennt die Rückfrage Wirkung und Umfang
    Und erst eine ausdrückliche Bestätigung setzt den Befehl ab

  @PG-10
  Szenario: Ein IP-Bann ist keine Kleinigkeit
    Angenommen ein Bann auf eine IP-Adresse
    Wenn er ausgelöst wird
    Dann weist die Rückfrage darauf hin, dass er mehr als eine Person treffen kann
    Und die Adresse steht sichtbar im Text

  @PG-10
  Szenario: Den Weltgenerator anzuhalten ist ein eigener Schritt
    Angenommen ein laufender Weltgenerator
    Wenn "stoppen" gewählt wird
    Dann fragt die Anwendung zurück
    Und der Zustand danach wird aus "worldgen status" gelesen, nicht angenommen

  @PG-11
  Szenario: Ein neuer Befehl spricht vier Sprachen
    Angenommen ein neu hinzugefügter Befehl mit Dialog und Meldungen
    Wenn die Sprachkataloge verglichen werden
    Dann existiert jeder neue Schlüssel in allen vier
    Und keiner davon steht als Klartext im Go-Code
