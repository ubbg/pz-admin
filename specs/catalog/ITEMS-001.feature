# language: de
# Spec: ITEMS-001 — Der Item-Katalog aus Build 42
# Jedes Szenario beschreibt das *Verhalten* einer Leitplanke und ist mit ihrer
# PG-ID getaggt. Die IDs hier müssen deckungsgleich mit `guardrails` in
# ITEMS-001.spec.md sein (Abgleich siehe specs/README.md).
#
# Diese Datei wird nicht ausgeführt — pz-admin hat keine Testsuite.

Funktionalität: Ein Katalog, der sagt, woher er kommt
  Der Item-Katalog wird gegen eine benannte Wiki-Fassung erzeugt und trägt diese
  Herkunft in sich. Fehlende Bilder halten den Lauf nicht auf, und wer eine ID
  kennt, die im Katalog fehlt, kommt trotzdem an sie heran.

  Grundlage:
    Angenommen der Generator "yarn items"
    Und eine Wiki-Fassung, die einem Spiel-Build entspricht

  @PG-12
  Szenario: Die erzeugte Datei nennt ihren Stand
    Angenommen ein abgeschlossener Katalog-Lauf
    Wenn die erzeugte Datei geöffnet wird
    Dann stehen Quell-Adresse, Fassung, Build und Erzeugungsdatum darin
    Und die Frage "welcher Stand?" ist ohne Blick in den Quelltext beantwortet

  @PG-12
  Szenario: Ein Lauf gegen die Live-Seite trägt trotzdem eine Fassung ein
    Angenommen ein Lauf gegen die laufend gepflegte Wiki-Seite
    Wenn der Katalog geschrieben wird
    Dann steht dort die tatsächlich ausgelieferte Fassungsnummer
    Und nicht das Wort "live"

  @PG-12
  Szenario: Build-42-Gegenstände sind enthalten
    Angenommen ein Katalog aus einer Build-42-Fassung
    Wenn nach Schmiede- und Tierhaltungsgegenständen gesucht wird
    Dann finden sich Amboss, Schmelztiegel und Sattel
    Und nicht nur der Bestand aus Build 41

  @PG-13
  Szenario: Ein fehlendes Bild hält den Lauf nicht auf
    Angenommen ein Katalog-Lauf, bei dem einzelne Bilder nicht ladbar sind
    Wenn der Lauf durchläuft
    Dann endet er ordentlich
    Und die Anzahl der fehlenden Bilder steht in der Schlussbilanz

  @PG-13
  Szenario: Die Bilanz macht den Schaden sichtbar
    Angenommen ein Lauf mit dreitausend Einträgen und zwölf fehlenden Bildern
    Wenn er endet
    Dann nennt die Ausgabe Einträge, geladene Bilder und Fehlende
    Und man erkennt ohne Diff, ob der Lauf brauchbar war

  @PG-14
  Szenario: Eine unbekannte Item-ID lässt sich trotzdem verschenken
    Angenommen ein Gegenstand aus einer Mod, der in keinem Wiki steht
    Wenn seine ID von Hand eingegeben wird
    Dann wird sie unverändert an den Server geschickt
    Und der Katalog steht dem nicht im Weg

  @PG-14
  Szenario: Groß- und Kleinschreibung bleibt erhalten
    Angenommen die eingegebene ID "Base.Axe"
    Wenn der Befehl abgesetzt wird
    Dann steht dort "Base.Axe"
    Und nicht die kleingeschriebene Form, die der Server ablehnt

  @PG-15
  Szenario: Kein Eintrag verweist auf ein Bild, das es nicht gibt
    Angenommen der erzeugte Katalog samt Bilderordner
    Wenn jeder Bildverweis geprüft wird
    Dann existiert zu jedem eine Datei
    Und die Oberfläche zeigt keine leere Fläche ohne Erklärung

  @PG-15
  Szenario: Alte Bilder aus früheren Läufen stören nicht
    Angenommen Bilddateien, auf die kein Eintrag mehr zeigt
    Wenn der Katalog geladen wird
    Dann ist das kein Fehler
    Und ein Aufräumen bleibt ein eigener, geprüfter Schritt
