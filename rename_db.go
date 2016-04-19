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
	"fmt"
)

func enumerateTables(tx *sql.Tx) ([]string, error) {
	var tables []string
	var field string
	rows, err := tx.Query("SHOW TABLES")
	if err != nil {
		return []string{}, fmt.Errorf("Could not show tables: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&field)
		if err != nil {
			return []string{}, fmt.Errorf("Could not scan existing tables: %v", err)
		}
		tables = append(tables, field)
	}
	err = rows.Err()
	if err != nil {
		return []string{}, fmt.Errorf("Could not read existing tables: %v", err)
	}
	return tables, nil
}

func renameDatabase(db *sql.DB, from, to string) (result error) {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(fmt.Sprintf("USE `%s`", from))
	if err != nil {
		return fmt.Errorf("could not select DB: %v", err)
	}

	tablesToMove, err := enumerateTables(tx)
	if err != nil {
		return fmt.Errorf("could not enumerate tables: %v", err)
	}

	// first attempt creation of the destination
	_, err = tx.Exec(fmt.Sprintf("CREATE DATABASE `%s`", to))
	if err != nil {
		return fmt.Errorf("could not create destination DB: %v", err)
	}

	// start moving tables
	for _, table := range tablesToMove {
		if _, err := tx.Exec(fmt.Sprintf("ALTER TABLE `%s`.`%s` RENAME `%s`.`%s`", from, table, to, table)); err != nil {
			return fmt.Errorf("could not move table '%s': %v", table, err)
		}
	}

	// drop original db - first make sure it's now empty
	tablesLeft, err := enumerateTables(tx)
	if err != nil {
		return fmt.Errorf("could not enumerate tables: %v", err)
	}

	if len(tablesLeft) != 0 {
		return fmt.Errorf("refusing to drop original database as it still contains tables: %v", err)
	}

	_, err = tx.Exec(fmt.Sprintf("DROP DATABASE `%s`", from))
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
