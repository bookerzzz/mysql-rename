/*
 * mysql-rename v0.1.0 - utility to rename MySQL databases
 * Copyright (C) 2016 bookerzzz - https://github.com/bookerzzz/mysql-rename

This program is free software; you can redistribute it and/or
modify it under the terms of the GNU General Public License
as published by the Free Software Foundation; either version 2
of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program; if not, write to the Free Software
Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.
*/
package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"

	"github.com/jessevdk/go-flags"
)

var cliFlags struct {
	MySQLDSN string `long:"mysql-dsn" env:"MYSQL_DSN" default:"root@tcp(localhost:3306)/" description:"MySQL Data Source Name URI"`
	From     string `long:"from" env:"FROM" default:"" description:"Database to be renamed"`
	To       string `long:"to" env:"TO" default:"" description:"New database name"`
}

func main() {
	var parser = flags.NewParser(&cliFlags, flags.Default)
	// parse flags
	args, err := parser.Parse()
	if err != nil {
		if et, ok := err.(*flags.Error); ok {
			if et.Type == flags.ErrHelp {
				return
			}
		}
		log.Fatalf("error parsing flags: %v", err)
		return
	}
	if len(args) > 0 {
		log.Fatalf("unexpected arguments: %v", args)
		return
	}

	if len(cliFlags.From) == 0 {
		log.Fatal("no source database name specified")
		return
	}
	if len(cliFlags.To) == 0 {
		log.Fatal("no destination database name specified")
		return
	}

	db, err := sql.Open("mysql", cliFlags.MySQLDSN)
	if err != nil {
		log.Fatalf("could not connect to MySQL: %v", err)
	}

	err = renameDatabase(db, cliFlags.From, cliFlags.To)
	if err != nil {
		log.Fatalf("could not rename '%s' to '%s': %v", cliFlags.From, cliFlags.To, err)
	}

	log.Printf("Renamed '%s' to '%s' successfully", cliFlags.From, cliFlags.To)
}
