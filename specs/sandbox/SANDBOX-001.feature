# language: de
# Spec: SANDBOX-001 — SandboxVars bearbeiten, was RCON nicht kann
# Jedes Szenario beschreibt das *Verhalten* einer Leitplanke und ist mit ihrer
# PG-ID getaggt. Die IDs hier müssen deckungsgleich mit `guardrails` in
# SANDBOX-001.spec.md sein (Abgleich siehe specs/README.md).
#
# Diese Datei wird nicht ausgeführt — pz-admin hat keine Testsuite.

Funktionalität: Die Sandbox-Datei bearbeiten, ohne sie zu beschädigen
  Sandbox-Einstellungen leben in einer Lua-Datei auf dem Server und nicht hinter
  RCON. Der Editor fasst sie nur an, wenn ihm gesagt wurde wo, sichert vorher,
  behält die Form und sagt, dass es erst nach dem Neustart gilt.

  Grundlage:
    Angenommen ein Server mit einer Datei "servertest_SandboxVars.lua"
    Und zweihundertdreiundsiebzig Variablen darin

  @PG-16
  Szenario: Sandbox-Werte gehen nicht über RCON
    Angenommen eine geänderte Sandbox-Variable
    Wenn gespeichert wird
    Dann wird kein "changeoption" an den Server geschickt
    Und die Änderung landet in der Datei

  @PG-16
  Szenario: Der freundliche Fehlschlag wird gar nicht erst versucht
    Angenommen eine Variable, die es früher in der ini gab
    Wenn sie im Sandbox-Editor geändert wird
    Dann nimmt die Anwendung den Dateiweg
    Und nicht den Weg über die Serverkonsole

  @PG-17
  Szenario: Vor dem Schreiben liegt eine Sicherung
    Angenommen ein Speichervorgang
    Wenn er beginnt
    Dann existiert neben der Datei eine Kopie mit Zeitstempel
    Und sie entspricht dem Stand vor der Änderung

  @PG-17
  Szenario: Ein abgebrochener Schreibvorgang ist kein Verlust
    Angenommen ein Speichervorgang, der auf halbem Weg scheitert
    Wenn der Fehler gemeldet wird
    Dann nennt die Meldung den Pfad der Sicherung
    Und der vorherige Stand ist wiederherstellbar

  @PG-18
  Szenario: Kommentare und Reihenfolge überleben
    Angenommen eine Datei mit Kommentaren zwischen den Werten
    Wenn genau eine Variable geändert und gespeichert wird
    Dann unterscheidet sich die Datei in genau dieser Zeile
    Und alle Kommentare stehen unverändert an ihrem Platz

  @PG-18
  Szenario: Untertabellen bleiben Untertabellen
    Angenommen die Untertabellen für Keller, Karte und Zombie-Verhalten
    Wenn ein Wert darin geändert wird
    Dann bleibt die Verschachtelung erhalten
    Und die Datei ist danach weiterhin von Hand lesbar

  @PG-19
  Szenario: Der Neustart-Hinweis steht dort, wo gespeichert wird
    Angenommen der Sandbox-Tab mit ungespeicherten Änderungen
    Wenn er angesehen wird
    Dann steht der Hinweis auf den nötigen Serverneustart sichtbar am Speichern-Knopf
    Und nicht in einer aufklappbaren Fußnote

  @PG-19
  Szenario: Nach dem Speichern wird nichts Falsches versprochen
    Angenommen ein erfolgreicher Speichervorgang
    Wenn die Bestätigung erscheint
    Dann sagt sie, dass die Änderung erst nach dem Neustart gilt
    Und nicht, dass sie übernommen wurde

  @PG-20
  Szenario: Ohne konfigurierten Zugang wird nichts geraten
    Angenommen eine Anwendung ohne eingerichteten Sandbox-Pfad
    Wenn der Sandbox-Tab geöffnet wird
    Dann ist er inaktiv und nennt den Grund
    Und es wurde keine Datei geöffnet und keine angelegt

  @PG-20
  Szenario: Der Zugangsweg wird ausdrücklich gewählt
    Angenommen die drei Wege lokaler Pfad, Docker-Volume und SFTP
    Wenn der Zugang eingerichtet wird
    Dann wählt die Nutzerin einen davon aus
    Und die Anwendung übernimmt keinen Standardpfad als Annahme

  @PG-20
  Szenario: Der Editor braucht keine RCON-Verbindung
    Angenommen eine getrennte RCON-Verbindung
    Und ein eingerichteter Sandbox-Zugang
    Wenn der Sandbox-Tab geöffnet wird
    Dann lässt sich die Datei trotzdem lesen und schreiben
    Und der Neustart-Hinweis gilt unverändert
