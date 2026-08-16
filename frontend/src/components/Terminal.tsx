import React, { useState, useRef, useEffect } from "react";
import { useRcon } from "@/contexts/rcon-provider";
import { ScrollArea } from "@/components/ui/scroll-area";
import { useConfig } from "@/contexts/config-provider";

// Befehlsliste für die Autovervollständigung.
//
// Gepflegter Stand: Project Zomboid Build 42 (pzwiki „Admin commands", Fassung
// 42.20.2) — 59 Befehle, dazu "cls" für das Leeren dieses Terminals. RCON gibt keine
// maschinenlesbare Befehlsliste her (/help liefert Fließtext), deshalb ist die Liste
// handgepflegt; sie wird gegen die Wiki-Seite abgeglichen, nicht geraten.
//
// Zwei Befehle schreibt das Wiki verkürzt (godmod, godmodplayer) und vermerkt dazu
// einen Schreibfehler im Spiel; beide Schreibweisen wirken. Hier steht die richtige.
const commands = [
  "additem",
  "addkey",
  "addsteamid",
  "addtosafehouse",
  "adduser",
  "addvehicle",
  "addxp",
  "alarm",
  "banid",
  "banip",
  "banuser",
  "changeoption",
  "checkModsNeedUpdate",
  "chopper",
  "cls",
  "createhorde",
  "createhorde2",
  "godmode",
  "godmodeplayer",
  "gunshot",
  "help",
  "invisible",
  "invisibleplayer",
  "kick",
  "kickfromsafehouse",
  "lightning",
  "list",
  "log",
  "noclip",
  "players",
  "quit",
  "releasesafehouse",
  "reloadalllua",
  "reloadlua",
  "reloadoptions",
  "remove",
  "removeitem",
  "removemapsymbolsforuser",
  "removesteamid",
  "removeuserfromwhitelist",
  "removezombies",
  "save",
  "servermsg",
  "setaccesslevel",
  "setpassword",
  "showoptions",
  "startrain",
  "startstorm",
  "stats",
  "stoprain",
  "stopweather",
  "teleport",
  "teleportplayer",
  "teleportto",
  "thunder",
  "unbanid",
  "unbanip",
  "unbanuser",
  "voiceban",
  "worldgen",
];

// Unterbefehle, die keine Spielernamen sind.
const worldgenSubcommands = ["start", "recheck", "stop", "status"];
const logLevels = ["Trace", "Debug", "General", "Warning", "Error"];
const statsModes = ["none", "file", "console", "all"];

const TerminalPage: React.FC = () => {
  const { config } = useConfig();
  const { sendCommand, players, optionNames } = useRcon();

  const playerNames = players.map((player) => player.name);

  const commandsMap: { [key: string]: string[] } = {
    empty: commands,
    help: commands,
    setaccesslevel: playerNames,
    banuser: playerNames,
    unbanuser: playerNames,
    kick: playerNames,
    kickfromsafehouse: playerNames,
    addtosafehouse: playerNames,
    teleport: playerNames,
    teleportplayer: playerNames,
    tp: playerNames,
    createhorde: playerNames,
    createhorde2: playerNames,
    lightning: playerNames,
    thunder: playerNames,
    addxp: playerNames,
    additem: playerNames,
    addkey: playerNames,
    addvehicle: playerNames,
    godmode: playerNames,
    godmodeplayer: playerNames,
    invisible: playerNames,
    invisibleplayer: playerNames,
    noclip: playerNames,
    voiceban: playerNames,
    setpassword: playerNames,
    removemapsymbolsforuser: playerNames,
    // Die Optionsnamen kommen vom verbundenen Server, nicht aus einer Liste im Code.
    changeoption: optionNames,
    worldgen: worldgenSubcommands,
    log: logLevels,
    stats: statsModes,
  };

  const [output, setOutput] = useState<{ type: "command" | "response" | "info" | "error"; line: string }[]>([
    { type: "info", line: "Welcome, type 'help' for commands or 'cls' to clear the terminal." },
    { type: "info", line: "You can use tab for auto-completion and up/down keys to navigate command history." },
  ]);
  const [currentInput, setCurrentInput] = useState<string>("");
  const [commandHistory, setCommandHistory] = useState<string[]>([]);
  const [historyIndex, setHistoryIndex] = useState<number>(-1);
  const [tabMatches, setTabMatches] = useState<string[]>([]);
  const [tabIndex, setTabIndex] = useState<number>(-1);

  const inputRef = useRef<HTMLInputElement>(null);
  const scrollRef = useRef<HTMLDivElement>(null);

  const scrollToBottom = () => {
    scrollRef.current?.scrollIntoView({ behavior: "smooth" });
  };

  useEffect(() => {
    inputRef.current?.focus();
    scrollToBottom();
  }, [output]);

  const addOutput = (line: string, type: "command" | "response" | "info" | "error") => {
    const newLines = line.split("\n").map((l) => ({ type, line: l }));
    setOutput((prev) => [...prev, ...newLines]);
  };

  const moveCommandToEnd = (command: string) => {
    setCommandHistory((prev) => {
      const filteredHistory = prev.filter((cmd) => cmd !== command);
      return [...filteredHistory, command];
    });
  };

  const handleTabCompletion = () => {
    const trimmedInput = currentInput.trim();
    const words = trimmedInput.split(" ");
    const baseCommand = words[0]; // The first word is the base command
    const currentWord = words[words.length - 1];

    const suggestions =
      words.length === 1
        ? commandsMap["empty"] // Suggest commands for the base command
        : commandsMap[baseCommand] || []; // Suggest values for a specific command

    if (tabIndex === -1) {
      const matches = suggestions.filter((item) => item.startsWith(currentWord));
      if (matches.length > 0) {
        setTabMatches(matches);
        setTabIndex(0);
        setCurrentInput((words.slice(0, -1).join(" ") + " " + matches[0]).trim());
      }
    } else {
      const newIndex = (tabIndex + 1) % tabMatches.length;
      setTabIndex(newIndex);
      setCurrentInput((words.slice(0, -1).join(" ") + " " + tabMatches[newIndex]).trim());
    }
  };

  const resetTabCompletion = (value: string) => {
    setCurrentInput(value);
    setTabMatches([]);
    setTabIndex(-1);
  };

  const handleKeyDown = async (e: React.KeyboardEvent) => {
    if (e.key === "Enter" && currentInput.trim()) {
      const command = currentInput.trim();
      moveCommandToEnd(command);
      addOutput(`> ${command}`, "command");
      resetTabCompletion("");
      setHistoryIndex(-1);

      try {
        if (command === "cls") {
          setOutput([]);
        } else {
          const { response, error } = await sendCommand(command);
          if (response && !error) {
            addOutput(response, "response");
          } else {
            addOutput(error || "No response.", "error");
          }
        }
      } catch (error) {
        addOutput(`Error: ${error}`, "error");
      }

      scrollToBottom();
    }

    if (e.key === "ArrowUp") {
      e.preventDefault();
      if (commandHistory.length > 0) {
        const newIndex = historyIndex === -1 ? commandHistory.length - 1 : Math.max(0, historyIndex - 1);
        setHistoryIndex(newIndex);
        resetTabCompletion(commandHistory[newIndex]);
      }
    } else if (e.key === "ArrowDown") {
      e.preventDefault();
      if (historyIndex > -1 && historyIndex < commandHistory.length - 1) {
        const newIndex = historyIndex + 1;
        setHistoryIndex(newIndex);
        resetTabCompletion(commandHistory[newIndex]);
      } else if (historyIndex === commandHistory.length - 1) {
        setHistoryIndex(-1);
        resetTabCompletion("");
      }
    }

    if (e.key === "Tab") {
      e.preventDefault();
      handleTabCompletion();
    }

    if (e.key === "Escape") {
      resetTabCompletion("");
      setHistoryIndex(-1);
    }
  };

  return (
    <div
      className={`w-[calc(100vw-16rem)] ${
        config?.useSystemTitleBar ? "h-[calc(100vh-3.5rem)]" : "h-[calc(100vh-5.5rem)]"
      } dark:bg-black/20 bg-white/20 font-mono p-2 pb-1`}
      onClick={() => inputRef.current?.focus()}
    >
      <ScrollArea className="h-[calc(100%-2rem)] overflow-auto pr-4">
        <div>
          {output.map((entry, index) => (
            <div
              key={index}
              className={`break-words whitespace-pre-wrap w-[calc(100vw-18rem)] ${
                entry.type === "command"
                  ? "dark:text-green-400 text-green-900 font-semibold"
                  : entry.type === "response"
                  ? "dark:text-green-700 text-green-600"
                  : entry.type === "info"
                  ? "text-blue-500"
                  : "text-red-500"
              }`}
            >
              {entry.line}
            </div>
          ))}
          <div ref={scrollRef}></div>
        </div>
      </ScrollArea>
      <div className="flex">
        <span className="dark:text-green-400 text-green-900 mr-2 h-[2rem] flex items-center font-semibold">$</span>
        <input
          ref={inputRef}
          type="text"
          spellCheck="false"
          value={currentInput}
          onChange={(e) => resetTabCompletion(e.target.value)}
          onKeyDown={handleKeyDown}
          className="w-full bg-transparent outline-none border-none dark:text-green-400 text-green-900 h-[2rem]"
          autoFocus
        />
      </div>
    </div>
  );
};

export default TerminalPage;
