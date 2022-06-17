import React from "react";
import Search from "./Search";
import SearchBar from "./SearchBar";

function Navbar() {
  return (
    <nav className="">
      <ul className="flex gap-x-6 text-slate-500">
        <li>
          <a href="/" className="nav-font nav-item-selected">
            Trades
          </a>
        </li>
        <li>
          <a href="/" className="nav-font">
            Insiders
          </a>
        </li>
        <li>
          <a href="/" className="nav-font">
            Companies
          </a>
        </li>
      </ul>
      {/* <Search /> */}
    </nav>
  );
}

export default Navbar;
