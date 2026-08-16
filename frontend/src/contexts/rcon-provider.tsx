import React, {
  createContext,
  useContext,
  useState,
  useCallback,
  ReactNode,
  Dispatch,
  SetStateAction,
  useEffect,
  useMemo,
} from "react";
import {
  ConnectRcon,
  DisconnectRcon,
  ExportOptionsDialog,
  ImportOptionsDialog,
  SendRconCommand,
  UpdateOptions,
} from "@/wailsjs/go/main/App";
import { main } from "@/wailsjs/go/models";
import { EventsOff, EventsOn } from "@/wailsjs/runtime/runtime";
import { deepEqual } from "@/lib/utils";
import { optionsMap } from "@/assets/options";
import {
  emptyServerOptions,
  OptionValue,
  optionValueFromString,
  optionValuesToStrings,
  serverOptionsFrom,
  ServerOptions,
} from "@/lib/options";

interface RconContextType {
  isConnected: boolean;
  isConnecting: boolean;
  setIsConnected: Dispatch<SetStateAction<boolean>>;
  connect: (credentials: main.Credentials) => Promise<boolean>;
  disconnect: () => Promise<boolean>;
  sendCommand: (command: string) => Promise<main.RconResponse>;
  ip: string;
  port: string;
  players: main.Player[];

  options: Record<string, OptionValue>;
  modifiedOptions: Record<string, OptionValue>;
  // Was der Server gemeldet hat: Namen in Serverreihenfolge und der Typ je Name.
  optionNames: string[];
  optionKinds: Record<string, string>;
  modifyOption: (key: string, value: OptionValue) => void;
  cancelModifiedOptions: () => void;
  optionsModified: boolean;
  updateOptions: (reload?: boolean) => Promise<boolean>;
  updatingOptions: boolean;
  optionsInvalid: boolean;

  reloadDoubleptions: () => void;
  reloadDoubleOptionsKey: number;

  importOptions: () => void;
  exportOptions: () => void;
}

const RconContext = createContext<RconContextType | undefined>(undefined);

export const RconProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
  const [isConnected, setIsConnected] = useState(false);
  const [isConnecting, setIsConnecting] = useState(false);
  const [ip, setIp] = useState("");
  const [port, setPort] = useState("");
  const [players, setPlayers] = useState<main.Player[]>([]);

  const [serverOptions, setServerOptions] = useState<ServerOptions>(emptyServerOptions);
  const [modifiedOptions, setModifiedOptions] = useState<Record<string, OptionValue>>({});
  const [updatingOptions, setUpdatingOptions] = useState(false);

  const [reloadDoubleOptionsKey, setReloadDoubleOptionsKey] = useState(0);

  const options = serverOptions.values;

  const optionsModified = useMemo(() => !deepEqual(options, modifiedOptions), [options, modifiedOptions]);
  const optionsInvalid = useMemo(() => {
    return Object.entries(modifiedOptions).some(
      ([key, value]) =>
        value === null ||
        value === undefined ||
        (typeof value === "number" &&
          (isNaN(value) ||
            value < (optionsMap.get(key)?.Range?.Min ?? -2147483648) ||
            value > (optionsMap.get(key)?.Range?.Max ?? 2147483647)))
    );
  }, [modifiedOptions]);

  useEffect(() => {
    const handleUpdatePlayers = (players: main.Player[]) => {
      setPlayers(players);
    };

    const handleUpdateOptions = (list: main.Option[]) => {
      const next = serverOptionsFrom(list);
      setServerOptions(next);
      setModifiedOptions(next.values);
    };

    EventsOn("update-players", handleUpdatePlayers);
    EventsOn("update-options-list", handleUpdateOptions);

    return () => {
      EventsOff("update-players");
      EventsOff("update-options-list");
    };
  }, []);

  const connect = useCallback(async (credentials: main.Credentials): Promise<boolean> => {
    try {
      setIsConnecting(true);
      const result = await ConnectRcon(credentials);
      if (result) {
        setIp(credentials.ip);
        setPort(credentials.port);
      }
      setIsConnected(result);
      setIsConnecting(false);
      return result;
    } catch (error) {
      console.error("Error connecting to RCON:", error);
      setIsConnecting(false);
      return false;
    }
  }, []);

  const disconnect = useCallback(async (): Promise<boolean> => {
    try {
      const result = await DisconnectRcon();
      setIsConnected(!result);
      setPlayers([]);
      setServerOptions(emptyServerOptions);
      setModifiedOptions({});
      setIp("");
      setPort("");
      setUpdatingOptions(false);
      return result;
    } catch (error) {
      console.error("Error disconnecting from RCON:", error);
      return false;
    }
  }, []);

  const sendCommand = useCallback(
    async (command: string): Promise<main.RconResponse> => {
      if (!isConnected) {
        console.warn("Cannot send command: Not connected to RCON");
        return {} as main.RconResponse;
      }
      try {
        return await SendRconCommand(command);
      } catch (error) {
        console.error("Error sending RCON command:", error);
        return {} as main.RconResponse;
      }
    },
    [isConnected]
  );

  const modifyOption = useCallback((key: string, value: OptionValue) => {
    setModifiedOptions((prevOptions) => ({ ...prevOptions, [key]: value }));
  }, []);

  const updateOptions: RconContextType["updateOptions"] = useCallback(
    async (reload) => {
      setUpdatingOptions(true);
      const success = await UpdateOptions(optionValuesToStrings(modifiedOptions, serverOptions.kinds), reload ?? false);
      setUpdatingOptions(false);

      return success;
    },
    [modifiedOptions, serverOptions]
  );

  const cancelModifiedOptions = useCallback(() => {
    setModifiedOptions(options);
  }, [options]);

  useEffect(() => {
    if (!isConnected) {
      setIp("");
      setPort("");
    }
  }, [isConnected]);

  const reloadDoubleptions = useCallback(() => {
    setReloadDoubleOptionsKey((prevKey) => prevKey + 1);
  }, [modifiedOptions]);

  const importOptions = () => {
    ImportOptionsDialog().then((response) => {
      if (response.success) {
        // Eingelesen wird nur, was der Server auch kennt — der Rest hätte kein Ziel.
        setModifiedOptions((prevOptions) => {
          const next = { ...prevOptions };
          for (const [name, value] of Object.entries(response.options)) {
            if (serverOptions.kinds[name] !== undefined) {
              next[name] = optionValueFromString(value, serverOptions.kinds[name]);
            }
          }
          return next;
        });
      }
    });
  };

  const exportOptions = () => {
    ExportOptionsDialog(optionValuesToStrings(modifiedOptions, serverOptions.kinds));
  };

  return (
    <RconContext.Provider
      value={{
        isConnected,
        isConnecting,
        setIsConnected,
        connect,
        disconnect,
        sendCommand,
        ip,
        port,
        players,
        options,
        modifiedOptions,
        optionNames: serverOptions.names,
        optionKinds: serverOptions.kinds,
        optionsModified,
        modifyOption,
        cancelModifiedOptions,
        updateOptions,
        updatingOptions,
        optionsInvalid,
        reloadDoubleptions,
        reloadDoubleOptionsKey,
        importOptions,
        exportOptions,
      }}
    >
      {children}
    </RconContext.Provider>
  );
};

export const useRcon = (): RconContextType => {
  const context = useContext(RconContext);
  if (!context) {
    throw new Error("useRcon must be used within a RconProvider");
  }
  return context;
};
