package classic

import (
	. "classicserver/classic/constants"
	"fmt"
	"log"
	"strings"

	"golang.org/x/exp/slices"
)

const COLOR_ESCAPE string = "&"
const (
	COLOR_BLACK      = "&0"
	COLOR_DARK_BLUE  = "&1"
	COLOR_DARK_GREEN = "&2"
	COLOR_DARK_TEAL  = "&3"
	COLOR_DARK_RED   = "&4"
	COLOR_PURPLE     = "&5"
	COLOR_GOLD       = "&6"
	COLOR_GRAY       = "&7"
	COLOR_DARK_GRAY  = "&8"
	COLOR_BLUE       = "&9"
	COLOR_GREEN      = "&a"
	COLOR_TEAL       = "&b"
	COLOR_RED        = "&c"
	COLOR_PINK       = "&d"
	COLOR_YELLOW     = "&e"
	COLOR_WHITE      = "&f"
)

func (server *ClassicServer) HandleChat(player *Player, message string) {
	if message[0] == '/' {
		message := message[1:]
		parts := strings.Split(message, " ")
		command := parts[0]
		args := parts[1:]
		server.HandleCommand(player, command, args)
	} else {
		formatted := fmt.Sprintf("%s%s%s: %s", COLOR_WHITE, player.Username, COLOR_GRAY, message)
		server.BroadcastMessage(player.Id, formatted)
	}
}

func (server *ClassicServer) HandleCommand(player *Player, command string, args []string) {

	switch strings.ToLower(command) {

	case "help":
		handleHelp(server, player, args)

	case "about":
		handleAbout(server, player, args)

	case "kick":
		if player.Mode != MODE_OP {
			player.SendMessage("%sThis command is only available to operators", COLOR_DARK_RED)
			return
		}
		handleKick(server, player, args)

	case "tp":
		if player.Mode != MODE_OP {
			player.SendMessage("%sThis command is only available to operators", COLOR_DARK_RED)
			return
		}
		handleTp(server, player, args)

	case "ban":
		if player.Mode != MODE_OP {
			player.SendMessage("%sThis command is only available to operators", COLOR_DARK_RED)
			return
		}
		handleBan(server, player, args)

	case "unban":
		if player.Mode != MODE_OP {
			player.SendMessage("%sThis command is only available to operators", COLOR_DARK_RED)
			return
		}
		handleUnban(server, player, args)

	case "op":
		if player.Mode != MODE_OP {
			player.SendMessage("%sThis command is only available to operators", COLOR_DARK_RED)
			return
		}
		handleOp(server, player, args)

	case "deop":
		if player.Mode != MODE_OP {
			player.SendMessage("%sThis command is only available to operators", COLOR_DARK_RED)
			return
		}
		handleDeop(server, player, args)

	case "setspawn":
		if player.Mode != MODE_OP {
			player.SendMessage("%sThis command is only available to operators", COLOR_DARK_RED)
			return
		}
		handleSetSpawn(server, player, args)

	case "saveworld":
		if player.Mode != MODE_OP {
			player.SendMessage("%sThis command is only available to operators", COLOR_DARK_RED)
			return
		}
		handleSaveWorld(server, player, args)

	default:
		player.SendMessage("%sUnknown Command \"%s\"", COLOR_RED, command)

	}

}

func handleHelp(server *ClassicServer, player *Player, args []string) {
	player.SendMessage("%sAvailable Commands:", COLOR_TEAL)
	player.SendMessage("%s - /help - Show available commands", COLOR_DARK_TEAL)
	player.SendMessage("%s - /about - Display server info", COLOR_DARK_TEAL)
	if player.Mode == MODE_OP {
		player.SendMessage("%sOperator Commands:", COLOR_TEAL)
		player.SendMessage("%s - /kick <username> [reason] - Disconnect user", COLOR_DARK_TEAL)
		player.SendMessage("%s - /ban <username> [reason] - Ban & Disconnect user", COLOR_DARK_TEAL)
		player.SendMessage("%s - /unban <username> - Unban user", COLOR_DARK_TEAL)
		player.SendMessage("%s - /tp <playerfrom> <playerto> - Teleport a player to another", COLOR_DARK_TEAL)
		player.SendMessage("%s - /op <username> - Make user an operator", COLOR_DARK_TEAL)
		player.SendMessage("%s - /deop <username> - Remove operator from user", COLOR_DARK_TEAL)
		player.SendMessage("%s - /setspawn - Set server world spawn, saving the world", COLOR_DARK_TEAL)
		player.SendMessage("%s - /saveworld - Save the world", COLOR_DARK_TEAL)
	}
}

func handleAbout(server *ClassicServer, player *Player, args []string) {
	player.SendMessage("%s%s", COLOR_TEAL, server.Settings.Name)
	player.SendMessage("%s%s", COLOR_DARK_TEAL, server.Settings.MOTD)
}

func handleKick(server *ClassicServer, player *Player, args []string) {
	if len(args) < 1 {
		player.SendMessage("%sInvalid command, Expected:", COLOR_RED)
		player.SendMessage("%s/kick <username> [reason]", COLOR_RED)
	} else {
		username := args[0]
		if username == player.Username {
			player.SendMessage("%sCannot kick self", COLOR_RED)
			return
		}

		reason := "You were kicked"
		if len(args) > 1 {
			reason = strings.Join(args[1:], " ")
		}

		target := server.GetPlayerFromName(username)
		if target != nil {
			target.Kick(reason)
			player.SendMessage("%sKicking %s", COLOR_GREEN, username)
			log.Printf("%s has Kicked %s\n", player.Username, target.Username)
		} else {
			player.SendMessage("%sCould not find player \"%s\"", COLOR_RED, username)
		}
	}
}

func handleBan(server *ClassicServer, player *Player, args []string) {
	if len(args) < 1 {
		player.SendMessage("%sInvalid command, Expected:", COLOR_RED)
		player.SendMessage("%s/ban <username> [reason]", COLOR_RED)
	} else {
		username := args[0]
		if slices.Contains(server.Bans, username) {
			player.SendMessage("%s%s is already banned", COLOR_RED, username)
			return
		}

		reason := "You were kicked"
		if len(args) > 1 {
			reason = strings.Join(args[1:], " ")
		}

		server.Ban(username, reason)
		player.SendMessage("%sBanning %s", COLOR_GREEN, username)
		log.Printf("%s has banned %s\n", player.Username, username)
	}
}

func handleUnban(server *ClassicServer, player *Player, args []string) {
	if len(args) < 1 {
		player.SendMessage("%sInvalid command, Expected:", COLOR_RED)
		player.SendMessage("%s/unban <username>", COLOR_RED)
	} else {
		username := args[0]
		if !slices.Contains(server.Bans, username) {
			player.SendMessage("%s%s is not banned", COLOR_RED, username)
			return
		}

		server.Unban(username)
		player.SendMessage("%sUnbanning %s", COLOR_GREEN, username)
		log.Printf("%s has unbanned %s\n", player.Username, username)
	}
}

func handleOp(server *ClassicServer, player *Player, args []string) {
	if len(args) < 1 {
		player.SendMessage("%sInvalid command, Expected:", COLOR_RED)
		player.SendMessage("%s/op <username>", COLOR_RED)
	} else {
		username := args[0]
		if slices.Contains(server.OPs, username) {
			player.SendMessage("%s%s is already an operator", COLOR_RED, username)
			return
		}

		server.AddOP(username)
		player.SendMessage("%sMaking %s an operator", COLOR_GREEN, username)
		log.Printf("%s has made %s an operator\n", player.Username, username)
	}
}

func handleDeop(server *ClassicServer, player *Player, args []string) {
	if len(args) < 1 {
		player.SendMessage("%sInvalid command, Expected:", COLOR_RED)
		player.SendMessage("%s/deop <username>", COLOR_RED)
	} else {
		username := args[0]
		if !slices.Contains(server.OPs, username) {
			player.SendMessage("%s%s is not an operator", COLOR_RED, username)
			return
		}

		server.RemoveOP(username)
		player.SendMessage("%sRemoving %s as an operator", COLOR_GREEN, username)
		log.Printf("%s has made removed %s as an operator\n", player.Username, username)
	}
}

func handleSetSpawn(server *ClassicServer, player *Player, args []string) {
	server.World.SpawnX = float64(player.X)
	server.World.SpawnY = float64(player.Y)
	server.World.SpawnZ = float64(player.Z)
	server.World.SpawnYaw = player.Yaw
	server.World.SpawnPitch = player.Pitch
	world := *server.World

	save := func() {
		if err := SaveWorld(&world); err != nil {
			player.SendMessage("%sWorld save failed", COLOR_RED)
			log.Println(err)
			log.Printf("Failed world save; attempted by %s via /setspawn\n", player.Username)
		} else {
			player.SendMessage("%sWorld spawn set", COLOR_TEAL)
			log.Printf("World spawn set to X:%f Y:%f Z:%f\n", player.X, player.Y, player.Z)
		}
	}

	go save()
}

func handleTp(server *ClassicServer, player *Player, args []string) {
	if len(args) < 2 {
		player.SendMessage("%sInvalid command, Expected:", COLOR_RED)
		player.SendMessage("%s/tp <playerfrom> <playerto>", COLOR_RED)
	} else {
		fromUsername := args[0]
		toUsername := args[1]

		from := server.GetPlayerFromName(fromUsername)
		to := server.GetPlayerFromName(toUsername)

		if from == nil {
			player.SendMessage("%sCould not find player \"%s\"", COLOR_RED, fromUsername)
			return
		}

		if to == nil {
			player.SendMessage("%sCould not find player \"%s\"", COLOR_RED, toUsername)
			return
		}

		from.Teleport(to.X, to.Y, to.Z, from.Yaw, from.Pitch)
		player.SendMessage("%sTeleporting %s to %s", COLOR_GREEN, fromUsername, toUsername)
	}
}

func handleSaveWorld(server *ClassicServer, player *Player, args []string) {
	player.SendMessage("%sSaving world...", COLOR_TEAL)
	world := *server.World

	save := func() {
		if err := SaveWorld(&world); err != nil {
			player.SendMessage("%sWorld save failed", COLOR_RED)
			log.Println(err)
			log.Printf("Failed world save; attempted by %s via /saveworld\n", player.Username)
		} else {
			player.SendMessage("%sWorld saved", COLOR_GREEN)
			log.Printf("World saved by %s\n", player.Username)
		}
	}

	go save()
}
