---
name: tdd
description: Use whenever writing or changing production code. Enforces red-green-refactor - a failing test first, then the minimal code to pass, then refactor. Invoke for any feature, bugfix, or behavior change.
---

# TDD: red, green, refactor

Follow this loop for every change. Never skip step 1.

1. RED. Write one small failing test for the next behavior. Run it. Show it fails for the right reason.
   In Mode A, stop here and let me write the implementation.
2. GREEN. Write the minimum code to pass that one test. No extra features. Run it. Show it green.
3. REFACTOR. Improve names and structure with tests still green. Run tests again.
4. COMMIT. One passing behavior is one commit (see skill: commit).

Rules:
- One behavior per cycle. Heavy setup is a design smell; tell me.
- Test behavior through the public interface, not private internals.
- Table-driven tests by default in Go. Name cases by the behavior they pin.
- Never edit a test to make failing code pass. If the test is wrong, say why and fix it as its own step.
- If you cannot write a failing test for a change, stop and explain why before writing any code.