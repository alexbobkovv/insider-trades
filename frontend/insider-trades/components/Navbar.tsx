import React from "react";

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
    </nav>
  );
}

export default Navbar;
