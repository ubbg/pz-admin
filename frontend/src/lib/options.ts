import { main } from "@/wailsjs/go/models";
import { formatWithMinimumOneDecimal } from "@/lib/utils";

// Der Server liefert jede Option als Text. Für die Bedienelemente wird daraus ein
// Wert des abgeleiteten Typs; beim Speichern geht derselbe Weg zurück.
export type OptionValue = boolean | number | string;

export type ServerOptions = {
  values: Record<string, OptionValue>;
  kinds: Record<string, string>;
  names: string[];
};

export const emptyServerOptions: ServerOptions = { values: {}, kinds: {}, names: [] };

export function optionValueFromString(value: string, kind: string): OptionValue {
  switch (kind) {
    case "Boolean":
      return value.toLowerCase() === "true";
    case "Integer":
      return parseInt(value, 10);
    case "Double":
      return parseFloat(value);
    default:
      return value;
  }
}

export function optionValueToString(value: OptionValue, kind: string): string {
  if (typeof value === "boolean") {
    return value ? "true" : "false";
  }
  if (typeof value === "number") {
    if (isNaN(value)) {
      return "";
    }
    return kind === "Double" ? formatWithMinimumOneDecimal(value) : String(value);
  }
  return value;
}

export function serverOptionsFrom(list: main.Option[]): ServerOptions {
  const values: Record<string, OptionValue> = {};
  const kinds: Record<string, string> = {};
  const names: string[] = [];

  for (const option of list) {
    values[option.name] = optionValueFromString(option.value, option.kind);
    kinds[option.name] = option.kind;
    names.push(option.name);
  }

  return { values, kinds, names };
}

export function optionValuesToStrings(
  values: Record<string, OptionValue>,
  kinds: Record<string, string>
): Record<string, string> {
  const result: Record<string, string> = {};
  for (const [name, value] of Object.entries(values)) {
    result[name] = optionValueToString(value, kinds[name] ?? "String");
  }
  return result;
}
