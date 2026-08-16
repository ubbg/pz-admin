import fs from "fs";
import path from "path";
import { promisify } from "util";
import iconv from "iconv-lite";

const readFile = promisify(fs.readFile);
const writeFile = promisify(fs.writeFile);

// Function to parse the translation text file
async function parseTranslationFile(sourcePath, targetPath, encoding = "UTF-8") {
  try {
    const fileBuffer = await readFile(sourcePath);
    const fileContent = iconv.decode(fileBuffer, encoding);

    // Build 42 liefert die Übersetzungen als JSON ({"Base.Axe": "Axt", …}), Build 41
    // als Lua-artige .txt mit ItemName_Base.Axe = "Axt". Beide Formen werden gelesen.
    if (sourcePath.endsWith(".json")) {
      const parsed = JSON.parse(fileContent);
      const count = Object.keys(parsed).length;
      await writeFile(targetPath, JSON.stringify(parsed, null, 2), "utf-8");
      console.log(`Successfully wrote ${count} translations to ${targetPath}`);
      return;
    }

    const regex = /ItemName_([\w.]+\.[\w.]+)\s*=\s*"(.*)"/g;
    const translations = {};

    let match;
    while ((match = regex.exec(fileContent)) !== null) {
      const key = match[1]; // e.g., "Base.223Box"
      const value = match[2]; // e.g., ".223 Kalibre Mermi Kutusu"
      translations[key] = value;
    }

    const jsonContent = JSON.stringify(translations, null, 2);
    await writeFile(targetPath, jsonContent, "utf-8");
    console.log(`Successfully wrote translations to ${targetPath}`);
  } catch (error) {
    console.error(`Error: ${error.message}`);
  }
}

async function main(locale, encoding) {
  if (!locale) {
    console.error("Please provide a locale parameter.");
    process.exit(1);
  }

  // Quelle ist die Übersetzungsdatei des Spiels, kopiert nach game-translations/:
  // Build 42 unter media/lua/shared/Translate/<SPRACHE>/ItemName.json, Build 41 als
  // ItemName_<SPRACHE>.txt. Die Spieldateien selbst gehören nicht ins Repo.
  const jsonSource = path.resolve(`./game-translations/${locale}.json`);
  const sourcePath = fs.existsSync(jsonSource)
    ? jsonSource
    : path.resolve(`./game-translations/${locale}.txt`);
  const targetPath = path.resolve(`./public/locales/${locale}/items.json`);

  fs.mkdirSync(path.dirname(targetPath), { recursive: true });

  await parseTranslationFile(sourcePath, targetPath, encoding);
}

const locale = process.argv[2];
const encoding = process.argv[3] || "UTF-8";
main(locale, encoding);
