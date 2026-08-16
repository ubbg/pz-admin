import { useEffect, useState } from "react";
import { GetUpdateSource } from "@/wailsjs/go/main/App";
import { main } from "@/wailsjs/go/models";

// Die Quelle ändert sich zur Laufzeit nicht — sie steckt als Konstante im Build.
// Deshalb wird sie einmal geholt und danach aus dem Modul bedient.
let cachedSource: Promise<main.UpdateSource> | undefined;

function fetchUpdateSource(): Promise<main.UpdateSource> {
  if (!cachedSource) {
    cachedSource = GetUpdateSource();
  }
  return cachedSource;
}

export function useUpdateSource(): main.UpdateSource | undefined {
  const [source, setSource] = useState<main.UpdateSource>();

  useEffect(() => {
    let active = true;
    fetchUpdateSource().then((value) => {
      if (active) {
        setSource(value);
      }
    });
    return () => {
      active = false;
    };
  }, []);

  return source;
}
