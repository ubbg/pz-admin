export type Option = {
  FieldName: string;
  Type:
    | "String"
    | "Text"
    | "Password"
    | "Integer"
    | "Double"
    | "Boolean"
    | "Information"
    | "ServerWelcomeMessage"
    | "Choice"
    | "MultipleChoice"
    | "SpawnItems";
  Default?: boolean | number | string;
  Range?: {
    Min: number;
    Max: number;
  };
  DisabledValue?: boolean | number | string;
  // Anzeigen, aber nicht ändern lassen — für Werte, deren Änderung die eigene
  // Verbindung kappt (RCON-Port und -Passwort).
  ReadOnly?: boolean;
  Requirements?: {
    FieldName: string;
    FieldValue: boolean | number | string;
  }[];
  Choices?: {
    Name: string;
    Value: boolean | number | string;
  }[];
};

export type Category = {
  name: string;
  options: Option[];
};

export type Options = {
  categories: Category[];
};

// Darstellungsangaben zu den Serveroptionen: Kategorie, Bedienelement, Wertebereich,
// Auswahlwerte, Abhängigkeiten. Diese Tabelle entscheidet *nicht*, welche Optionen es
// gibt — das sagt der Server (siehe rcon_options.go). Fehlt hier ein Eintrag,
// erscheint die Option unter „Other Options" mit ihrem Rohnamen.
//
// Stand: Project Zomboid Build 42 (pzwiki „Server settings", Fassung 42.20.0,
// 144 Optionen).
export const options: Options = {
  categories: [
    {
      name: "General",
      options: [
        {
          FieldName: "PublicName",
          Type: "String",
        },
        {
          FieldName: "PublicDescription",
          Type: "Text",
        },
        {
          FieldName: "ServerWelcomeMessage",
          Type: "ServerWelcomeMessage",
          Default:
            "Welcome to Project Zomboid Multiplayer! <LINE> <LINE> To interact with the Chat panel: press Tab, T, or Enter. <LINE> <LINE> The Tab key will change the target stream of the message. <LINE> <LINE> Global Streams: /all <LINE> Local Streams: /say, /yell <LINE> Special Steams: /whisper, /safehouse, /faction. <LINE> <LINE> Press the Up arrow to cycle through your message history. Click the Gear icon to customize chat. <LINE> <LINE> Happy surviving!",
        },
        {
          FieldName: "Open",
          Type: "Boolean",
          Default: true,
        },
        {
          FieldName: "Public",
          Type: "Boolean",
          Default: false,
        },
        {
          FieldName: "Password",
          Type: "Password",
        },
        {
          FieldName: "DenyLoginOnOverloadedServer",
          Type: "Boolean",
          Default: true,
        },
        {
          FieldName: "SaveWorldEveryMinutes",
          Type: "Integer",
          Default: 0,
          DisabledValue: 0,
          Range: {
            Min: 0,
            Max: 2147483647,
          },
        },
        {
          FieldName: "Seed",
          Type: "String",
        },
        {
          FieldName: "SpawnItems",
          Type: "SpawnItems",
        },
        {
          FieldName: "SpawnPoint",
          Type: "String",
          Default: "0,0,0",
          DisabledValue: "0,0,0",
        },
        {
          FieldName: "MaxPlayers",
          Type: "Integer",
          Default: 32,
          Range: {
            Min: 1,
            Max: 100,
          },
        },
        {
          FieldName: "PauseEmpty",
          Type: "Boolean",
          Default: true,
        },
        {
          FieldName: "AllowCoop",
          Type: "Boolean",
          Default: true,
        },
        {
          FieldName: "AllowNonAsciiUsername",
          Type: "Boolean",
          Default: false,
        },
        {
          FieldName: "AnnounceDeath",
          Type: "Boolean",
          Default: false,
        },
        {
          FieldName: "AnnounceAnimalDeath",
          Type: "Boolean",
          Default: false,
        },
        {
          FieldName: "BanKickGlobalSound",
          Type: "Boolean",
          Default: true,
        },
        {
          FieldName: "Mods",
          Type: "String",
        },
        {
          FieldName: "WorkshopItems",
          Type: "String",
        },
        {
          FieldName: "Map",
          Type: "String",
          Default: "Muldraugh, KY",
        },
        {
          FieldName: "ClientCommandFilter",
          Type: "String",
          Default: "-vehicle.*;+vehicle.damageWindow;+vehicle.fixPart;+vehicle.installPart;+vehicle.uninstallPart",
        },
        {
          FieldName: "ClientActionLogs",
          Type: "String",
          Default: "ISEnterVehicle;ISExitVehicle;ISTakeEngineParts;",
        },
        {
          FieldName: "PerkLogs",
          Type: "Boolean",
          Default: true,
        },
      ],
    },
    {
      name: "Gameplay & Mechanics",
      options: [
        {
          FieldName: "CarEngineAttractionModifier",
          Type: "Double",
          Default: 0.5,
          Range: {
            Min: 0,
            Max: 10,
          },
        },
        {
          FieldName: "SpeedLimit",
          Type: "Double",
          Default: 70,
          Range: {
            Min: 10,
            Max: 150,
          },
        },
        {
          FieldName: "ItemNumbersLimitPerContainer",
          Type: "Integer",
          Default: 0,
          DisabledValue: 0,
          Range: {
            Min: 0,
            Max: 9000,
          },
        },
        {
          FieldName: "AllowDestructionBySledgehammer",
          Type: "Boolean",
          Default: true,
        },
        {
          FieldName: "SledgehammerOnlyInSafehouse",
          Type: "Boolean",
          Default: false,
          Requirements: [
            {
              FieldName: "AllowDestructionBySledgehammer",
              FieldValue: true,
            },
          ],
        },
        {
          FieldName: "DisableVehicleTowing",
          Type: "Boolean",
          Default: false,
        },
        {
          FieldName: "DisableTrailerTowing",
          Type: "Boolean",
          Default: false,
        },
        {
          FieldName: "DisableBurntTowing",
          Type: "Boolean",
          Default: false,
        },
        {
          FieldName: "NoFire",
          Type: "Boolean",
          Default: false,
        },
        {
          FieldName: "BloodSplatLifespanDays",
          Type: "Integer",
          Default: 0,
          DisabledValue: 0,
          Range: {
            Min: 0,
            Max: 365,
          },
        },
        {
          FieldName: "SleepAllowed",
          Type: "Boolean",
          Default: false,
        },
        {
          FieldName: "SleepNeeded",
          Type: "Boolean",
          Default: false,
          Requirements: [
            {
              FieldName: "SleepAllowed",
              FieldValue: true,
            },
          ],
        },
        {
          FieldName: "FastForwardMultiplier",
          Type: "Double",
          Default: 40,
          Range: {
            Min: 1,
            Max: 100,
          },
          Requirements: [
            {
              FieldName: "SleepAllowed",
              FieldValue: true,
            },
          ],
        },
        {
          FieldName: "UltraSpeedDoesnotAffectToAnimals",
          Type: "Boolean",
          Default: false,
        },
        {
          FieldName: "MapRemotePlayerVisibility",
          Type: "Choice",
          Default: 1,
          Choices: [
            {
              Name: "Hidden",
              Value: 1,
            },
            {
              Name: "Friends",
              Value: 2,
            },
            {
              Name: "Everyone",
              Value: 3,
            },
          ],
        },
        {
          FieldName: "HidePlayersBehindYou",
          Type: "Boolean",
          Default: true,
        },
        {
          FieldName: "PlayerBumpPlayer",
          Type: "Boolean",
          Default: false,
        },
        {
          FieldName: "KnockedDownAllowed",
          Type: "Boolean",
          Default: true,
        },
        {
          FieldName: "SneakModeHideFromOtherPlayers",
          Type: "Boolean",
          Default: true,
        },
        {
          FieldName: "PlayerRespawnWithOther",
          Type: "Boolean",
          Default: false,
          Requirements: [
            {
              FieldName: "AllowCoop",
              FieldValue: true,
            },
          ],
        },
        {
          FieldName: "PlayerRespawnWithSelf",
          Type: "Boolean",
          Default: false,
        },
        {
          FieldName: "RemovePlayerCorpsesOnCorpseRemoval",
          Type: "Boolean",
          Default: false,
        },
        {
          FieldName: "TrashDeleteAll",
          Type: "Boolean",
          Default: false,
        },
        {
          FieldName: "UsePhysicsHitReaction",
          Type: "Boolean",
          Default: false,
        },
        {
          FieldName: "SwitchZombiesOwnershipEachUpdate",
          Type: "Boolean",
          Default: false,
        },
      ],
    },
    {
      name: "Safehouse",
      options: [
        {
          FieldName: "PlayerSafehouse",
          Type: "Boolean",
          Default: false,
        },
        {
          FieldName: "AdminSafehouse",
          Type: "Boolean",
          Default: false,
          Requirements: [
            {
              FieldName: "PlayerSafehouse",
              FieldValue: false,
            },
          ],
        },
        {
          FieldName: "SafehouseDaySurvivedToClaim",
          Type: "Integer",
          Default: 0,
          Range: {
            Min: 0,
            Max: 2147483647,
          },
        },
        {
          FieldName: "SafeHouseRemovalTime",
          Type: "Integer",
          Default: 144,
          Range: {
            Min: 0,
            Max: 2147483647,
          },
        },
        {
          FieldName: "DisableSafehouseWhenOwnerConnected",
          Type: "Boolean",
          Default: false,
          Requirements: [
            {
              FieldName: "PlayerSafehouse",
              FieldValue: true,
            },
          ],
        },
        {
          FieldName: "SafehouseAllowNonResidential",
          Type: "Boolean",
          Default: false,
        },
        {
          FieldName: "SafehouseAllowRespawn",
          Type: "Boolean",
          Default: false,
        },
        {
          FieldName: "SafehouseAllowFire",
          Type: "Boolean",
          Default: true,
        },
        {
          FieldName: "SafehouseAllowTrepass",
          Type: "Boolean",
          Default: true,
        },
        {
          FieldName: "SafehouseAllowLoot",
          Type: "Boolean",
          Default: true,
        },
        {
          FieldName: "SafehousePreventsLootRespawn",
          Type: "Boolean",
          Default: true,
        },
        {
          FieldName: "SafehouseDisableDisguises",
          Type: "Boolean",
          Default: true,
          Requirements: [
            {
              FieldName: "PlayerSafehouse",
              FieldValue: true,
            },
          ],
        },
        {
          FieldName: "MaxSafezoneSize",
          Type: "Integer",
          Default: 20000,
          Range: {
            Min: 0,
            Max: 2147483647,
          },
        },
      ],
    },
    {
      name: "War",
      options: [
        {
          FieldName: "War",
          Type: "Boolean",
          Default: false,
        },
        {
          FieldName: "WarStartDelay",
          Type: "Integer",
          Default: 600,
          Range: {
            Min: 60,
            Max: 2147483647,
          },
          Requirements: [
            {
              FieldName: "War",
              FieldValue: true,
            },
          ],
        },
        {
          FieldName: "WarDuration",
          Type: "Integer",
          Default: 3600,
          Range: {
            Min: 60,
            Max: 2147483647,
          },
          Requirements: [
            {
              FieldName: "War",
              FieldValue: true,
            },
          ],
        },
        {
          FieldName: "WarSafehouseHitPoints",
          Type: "Integer",
          Default: 3,
          Range: {
            Min: 0,
            Max: 2147483647,
          },
          Requirements: [
            {
              FieldName: "War",
              FieldValue: true,
            },
          ],
        },
      ],
    },
    {
      name: "Faction",
      options: [
        {
          FieldName: "Faction",
          Type: "Boolean",
          Default: true,
        },
        {
          FieldName: "FactionDaySurvivedToCreate",
          Type: "Integer",
          Default: 0,
          Requirements: [
            {
              FieldName: "Faction",
              FieldValue: true,
            },
          ],
          Range: {
            Min: 0,
            Max: 2147483647,
          },
        },
        {
          FieldName: "FactionPlayersRequiredForTag",
          Type: "Integer",
          Default: 1,
          Requirements: [
            {
              FieldName: "Faction",
              FieldValue: true,
            },
          ],
          Range: {
            Min: 1,
            Max: 2147483647,
          },
        },
      ],
    },
    {
      name: "Player",
      options: [
        {
          FieldName: "DisplayUserName",
          Type: "Boolean",
          Default: true,
          Requirements: [
            {
              FieldName: "ShowFirstAndLastName",
              FieldValue: false,
            },
          ],
        },
        {
          FieldName: "ShowFirstAndLastName",
          Type: "Boolean",
          Default: false,
        },
        {
          FieldName: "MouseOverToSeeDisplayName",
          Type: "Boolean",
          Default: true,
        },
        {
          FieldName: "UsernameDisguises",
          Type: "Boolean",
          Default: false,
        },
        {
          FieldName: "HideDisguisedUserName",
          Type: "Boolean",
          Default: false,
          Requirements: [
            {
              FieldName: "UsernameDisguises",
              FieldValue: true,
            },
          ],
        },
        {
          FieldName: "HideAdminsInPlayerList",
          Type: "Boolean",
          Default: false,
        },
        {
          FieldName: "ShowCoordinates",
          Type: "Boolean",
          Default: false,
        },
        {
          FieldName: "DisableScoreboard",
          Type: "Boolean",
          Default: false,
        },
        {
          FieldName: "LoginQueueEnabled",
          Type: "Boolean",
          Default: false,
        },
        {
          FieldName: "LoginQueueConnectTimeout",
          Type: "Integer",
          Default: 60,
          Requirements: [
            {
              FieldName: "LoginQueueEnabled",
              FieldValue: true,
            },
          ],
          Range: {
            Min: 20,
            Max: 1200,
          },
        },
        {
          FieldName: "DropOffWhiteListAfterDeath",
          Type: "Boolean",
          Default: false,
        },
        {
          FieldName: "MaxAccountsPerUser",
          Type: "Integer",
          Default: 0,
          DisabledValue: 0,
          Range: {
            Min: 0,
            Max: 2147483647,
          },
        },
        {
          FieldName: "PingLimit",
          Type: "Integer",
          Default: 400,
          DisabledValue: 100,
          Range: {
            Min: 100,
            Max: 2147483647,
          },
        },
        {
          FieldName: "SteamScoreboard",
          Type: "Boolean",
          Default: true,
        },
      ],
    },
    {
      name: "PVP",
      options: [
        {
          FieldName: "PVP",
          Type: "Boolean",
          Default: true,
        },
        {
          FieldName: "SafetySystem",
          Type: "Boolean",
          Default: true,
          Requirements: [
            {
              FieldName: "PVP",
              FieldValue: true,
            },
          ],
        },
        {
          FieldName: "ShowSafety",
          Type: "Boolean",
          Default: true,
          Requirements: [
            {
              FieldName: "PVP",
              FieldValue: true,
            },
            {
              FieldName: "SafetySystem",
              FieldValue: true,
            },
          ],
        },
        {
          FieldName: "SafetyCooldownTimer",
          Type: "Integer",
          Default: 3,
          Requirements: [
            {
              FieldName: "PVP",
              FieldValue: true,
            },
            {
              FieldName: "SafetySystem",
              FieldValue: true,
            },
          ],
          Range: {
            Min: 0,
            Max: 1000,
          },
        },
        {
          FieldName: "SafetyToggleTimer",
          Type: "Integer",
          Default: 2,
          Requirements: [
            {
              FieldName: "PVP",
              FieldValue: true,
            },
            {
              FieldName: "SafetySystem",
              FieldValue: true,
            },
          ],
          Range: {
            Min: 0,
            Max: 1000,
          },
        },
        {
          FieldName: "SafetyDisconnectDelay",
          Type: "Integer",
          Default: 60,
          Range: {
            Min: 0,
            Max: 60,
          },
          Requirements: [
            {
              FieldName: "PVP",
              FieldValue: true,
            },
            {
              FieldName: "SafetySystem",
              FieldValue: true,
            },
          ],
        },
        {
          FieldName: "PVPFirearmDamageModifier",
          Type: "Double",
          Default: 50,
          Requirements: [
            {
              FieldName: "PVP",
              FieldValue: true,
            },
          ],
          Range: {
            Min: 0,
            Max: 500,
          },
        },
        {
          FieldName: "PVPMeleeDamageModifier",
          Type: "Double",
          Default: 30,
          Requirements: [
            {
              FieldName: "PVP",
              FieldValue: true,
            },
          ],
          Range: {
            Min: 0,
            Max: 500,
          },
        },
        {
          FieldName: "PVPMeleeWhileHitReaction",
          Type: "Boolean",
          Default: false,
          Requirements: [
            {
              FieldName: "PVP",
              FieldValue: true,
            },
          ],
        },
        {
          FieldName: "PVPLogToolChat",
          Type: "Boolean",
          Default: true,
          Requirements: [
            {
              FieldName: "PVP",
              FieldValue: true,
            },
          ],
        },
        {
          FieldName: "PVPLogToolFile",
          Type: "Boolean",
          Default: true,
          Requirements: [
            {
              FieldName: "PVP",
              FieldValue: true,
            },
          ],
        },
      ],
    },
    {
      name: "VOIP & Chat",
      options: [
        {
          FieldName: "GlobalChat",
          Type: "Boolean",
          Default: true,
        },
        {
          FieldName: "ChatStreams",
          Type: "MultipleChoice",
          Default: "s,r,w,y,sh,f,all",
          Choices: [
            {
              Name: "Say",
              Value: "s",
            },
            {
              Name: "Radio",
              Value: "r",
            },
            {
              Name: "Whisper",
              Value: "w",
            },
            {
              Name: "Yell",
              Value: "y",
            },
            {
              Name: "Safehouse",
              Value: "sh",
            },
            {
              Name: "Faction",
              Value: "f",
            },
            {
              Name: "All",
              Value: "all",
            },
          ],
        },
        {
          FieldName: "ChatMessageCharacterLimit",
          Type: "Integer",
          Default: 200,
          Range: {
            Min: 64,
            Max: 1024,
          },
        },
        {
          FieldName: "ChatMessageSlowModeTime",
          Type: "Integer",
          Default: 3,
          Range: {
            Min: 1,
            Max: 30,
          },
        },
        {
          FieldName: "BadWordListFile",
          Type: "String",
        },
        {
          FieldName: "GoodWordListFile",
          Type: "String",
        },
        {
          FieldName: "BadWordPolicy",
          Type: "Choice",
          Default: 3,
          Choices: [
            {
              Name: "ban",
              Value: 1,
            },
            {
              Name: "kick",
              Value: 2,
            },
            {
              Name: "log",
              Value: 3,
            },
            {
              Name: "mute",
              Value: 4,
            },
          ],
        },
        {
          FieldName: "BadWordReplacement",
          Type: "String",
          Default: "[HIDDEN]",
        },
        {
          FieldName: "DisableRadioInvisible",
          Type: "Boolean",
          Default: true,
        },
        {
          FieldName: "DisableRadioStaff",
          Type: "Boolean",
          Default: false,
        },
        {
          FieldName: "DisableRadioAdmin",
          Type: "Boolean",
          Default: true,
          Requirements: [
            {
              FieldName: "DisableRadioStaff",
              FieldValue: false,
            },
          ],
        },
        {
          FieldName: "DisableRadioModerator",
          Type: "Boolean",
          Default: false,
          Requirements: [
            {
              FieldName: "DisableRadioStaff",
              FieldValue: false,
            },
          ],
        },
        {
          FieldName: "DisableRadioOverseer",
          Type: "Boolean",
          Default: false,
          Requirements: [
            {
              FieldName: "DisableRadioStaff",
              FieldValue: false,
            },
          ],
        },
        {
          FieldName: "DisableRadioGM",
          Type: "Boolean",
          Default: true,
          Requirements: [
            {
              FieldName: "DisableRadioStaff",
              FieldValue: false,
            },
          ],
        },
        {
          FieldName: "VoiceEnable",
          Type: "Boolean",
          Default: true,
        },
        {
          FieldName: "Voice3D",
          Type: "Boolean",
          Default: true,
          Requirements: [
            {
              FieldName: "VoiceEnable",
              FieldValue: true,
            },
          ],
        },
        {
          FieldName: "VoiceMinDistance",
          Type: "Double",
          Default: 10,
          Requirements: [
            {
              FieldName: "VoiceEnable",
              FieldValue: true,
            },
          ],
          Range: {
            Min: 10,
            Max: 100000,
          },
        },
        {
          FieldName: "VoiceMaxDistance",
          Type: "Double",
          Default: 100,
          Requirements: [
            {
              FieldName: "VoiceEnable",
              FieldValue: true,
            },
          ],
          Range: {
            Min: 100,
            Max: 100000,
          },
        },
      ],
    },
    {
      name: "Discord",
      options: [
        {
          FieldName: "DiscordEnable",
          Type: "Boolean",
          Default: false,
          DisabledValue: "",
        },
        {
          FieldName: "DiscordToken",
          Type: "String",
          Requirements: [
            {
              FieldName: "DiscordEnable",
              FieldValue: true,
            },
          ],
        },
        {
          FieldName: "DiscordChatChannel",
          Type: "String",
          Requirements: [
            {
              FieldName: "DiscordEnable",
              FieldValue: true,
            },
          ],
        },
        {
          FieldName: "DiscordLogChannel",
          Type: "String",
          Requirements: [
            {
              FieldName: "DiscordEnable",
              FieldValue: true,
            },
          ],
        },
        {
          FieldName: "DiscordCommandChannel",
          Type: "String",
          Requirements: [
            {
              FieldName: "DiscordEnable",
              FieldValue: true,
            },
          ],
        },
        {
          FieldName: "WebhookAddress",
          Type: "String",
        },
      ],
    },
    {
      name: "Backup",
      options: [
        {
          FieldName: "BackupsOnStart",
          Type: "Boolean",
          Default: true,
        },
        {
          FieldName: "BackupsOnVersionChange",
          Type: "Boolean",
          Default: true,
        },
        {
          FieldName: "BackupsPeriod",
          Type: "Integer",
          Default: 0,
          DisabledValue: 0,
          Range: {
            Min: 0,
            Max: 1500,
          },
        },
        {
          FieldName: "BackupsCount",
          Type: "Integer",
          Default: 5,
          Range: {
            Min: 1,
            Max: 300,
          },
        },
      ],
    },
    {
      name: "Anti-Cheat",
      options: [
        {
          FieldName: "SteamVAC",
          Type: "Boolean",
          Default: true,
        },
        {
          FieldName: "DoLuaChecksum",
          Type: "Boolean",
          Default: true,
        },
        {
          FieldName: "AntiCheatSafety",
          Type: "Choice",
          Default: 2,
          Choices: [
            {
              Name: "ban",
              Value: 1,
            },
            {
              Name: "kick",
              Value: 2,
            },
            {
              Name: "log",
              Value: 3,
            },
            {
              Name: "disable",
              Value: 4,
            },
          ],
        },
        {
          FieldName: "AntiCheatSpeed",
          Type: "Choice",
          Default: 2,
          Choices: [
            {
              Name: "ban",
              Value: 1,
            },
            {
              Name: "kick",
              Value: 2,
            },
            {
              Name: "log",
              Value: 3,
            },
            {
              Name: "disable",
              Value: 4,
            },
          ],
        },
        {
          FieldName: "AntiCheatNoClip",
          Type: "Choice",
          Default: 4,
          Choices: [
            {
              Name: "ban",
              Value: 1,
            },
            {
              Name: "kick",
              Value: 2,
            },
            {
              Name: "log",
              Value: 3,
            },
            {
              Name: "disable",
              Value: 4,
            },
          ],
        },
        {
          FieldName: "AntiCheatHit",
          Type: "Choice",
          Default: 2,
          Choices: [
            {
              Name: "ban",
              Value: 1,
            },
            {
              Name: "kick",
              Value: 2,
            },
            {
              Name: "log",
              Value: 3,
            },
            {
              Name: "disable",
              Value: 4,
            },
          ],
        },
        {
          FieldName: "AntiCheatPacketException",
          Type: "Choice",
          Default: 4,
          Choices: [
            {
              Name: "ban",
              Value: 1,
            },
            {
              Name: "kick",
              Value: 2,
            },
            {
              Name: "log",
              Value: 3,
            },
            {
              Name: "disable",
              Value: 4,
            },
          ],
        },
        {
          FieldName: "AntiCheatPermission",
          Type: "Choice",
          Default: 2,
          Choices: [
            {
              Name: "ban",
              Value: 1,
            },
            {
              Name: "kick",
              Value: 2,
            },
            {
              Name: "log",
              Value: 3,
            },
            {
              Name: "disable",
              Value: 4,
            },
          ],
        },
        {
          FieldName: "AntiCheatXP",
          Type: "Choice",
          Default: 2,
          Choices: [
            {
              Name: "ban",
              Value: 1,
            },
            {
              Name: "kick",
              Value: 2,
            },
            {
              Name: "log",
              Value: 3,
            },
            {
              Name: "disable",
              Value: 4,
            },
          ],
        },
        {
          FieldName: "AntiCheatSafeHouse",
          Type: "Choice",
          Default: 2,
          Choices: [
            {
              Name: "ban",
              Value: 1,
            },
            {
              Name: "kick",
              Value: 2,
            },
            {
              Name: "log",
              Value: 3,
            },
            {
              Name: "disable",
              Value: 4,
            },
          ],
        },
        {
          FieldName: "AntiCheatPlayer",
          Type: "Choice",
          Default: 2,
          Choices: [
            {
              Name: "ban",
              Value: 1,
            },
            {
              Name: "kick",
              Value: 2,
            },
            {
              Name: "log",
              Value: 3,
            },
            {
              Name: "disable",
              Value: 4,
            },
          ],
        },
        {
          FieldName: "AntiCheatChecksum",
          Type: "Choice",
          Default: 2,
          Choices: [
            {
              Name: "ban",
              Value: 1,
            },
            {
              Name: "kick",
              Value: 2,
            },
            {
              Name: "log",
              Value: 3,
            },
            {
              Name: "disable",
              Value: 4,
            },
          ],
        },
      ],
    },
    {
      name: "Network",
      options: [
        {
          FieldName: "DefaultPort",
          Type: "Integer",
          Default: 16261,
          Range: {
            Min: 0,
            Max: 65535,
          },
        },
        {
          FieldName: "UDPPort",
          Type: "Integer",
          Default: 16262,
          Range: {
            Min: 0,
            Max: 65535,
          },
        },
        {
          FieldName: "UPnP",
          Type: "Boolean",
          Default: false,
        },
        {
          FieldName: "server_browser_announced_ip",
          Type: "String",
        },
        {
          FieldName: "RCONPort",
          Type: "Integer",
          Default: 27015,
          Range: {
            Min: 0,
            Max: 65535,
          },
          ReadOnly: true,
        },
        {
          FieldName: "RCONPassword",
          Type: "Password",
          ReadOnly: true,
        },
        {
          FieldName: "MaxPacketsPerSecond",
          Type: "Integer",
          Default: 300,
          Range: {
            Min: 100,
            Max: 1000,
          },
        },
        {
          FieldName: "MultiplayerStatisticsPeriod",
          Type: "Integer",
          Default: 1,
          Range: {
            Min: 0,
            Max: 10,
          },
          DisabledValue: 0,
        },
      ],
    },
    {
      name: "Miscellaneous",
      options: [
        {
          FieldName: "ResetID",
          Type: "Information",
        },
        {
          FieldName: "ServerPlayerID",
          Type: "Information",
        },
      ],
    },
  ],
};

export const optionsFlat: Option[] = options.categories.flatMap((category) => category.options);
export const optionsMap: Map<string, Option> = new Map(optionsFlat.map((option) => [option.FieldName, option]));
