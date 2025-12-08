// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

/// @title Counter
/// @notice A simple counter contract deployed on zkSync L2. It exposes a single
/// integer value that can be incremented and queried. This minimal example is
/// intentionally straightforward to illustrate how to interact with zkSync from
/// off‑chain code.
contract Counter {
    /// @dev Storage for the current counter value.
    uint256 public value;

    /// @notice Increments the counter by one. Anyone can call this method.
    function increment() external {
        unchecked {
            value += 1;
        }
    }

    /// @notice Returns the current counter value.
    /// @return The value of the counter.
    function get() external view returns (uint256) {
        return value;
    }
}