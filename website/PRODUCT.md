# DBStudio Lite website

<!-- impeccable:product-schema 1 -->

## Platform

web

## Product Purpose

The website introduces DBStudio Lite and provides desktop release downloads. DBStudio Lite is a lightweight database client for PostgreSQL, MySQL, and SQLite, built with Go, Next.js, and Wails. The website is separate from the desktop frontend in `../frontend`.

## Users

Developers exploring database tables, inspecting schemas, and writing SQL. Audience inferred from the repository README and existing landing-page copy.

## Capabilities and Constraints

Repository-supported capabilities: local connection profiles, schema and column inspection, tabular data browsing, an offline-bundled Monaco SQL editor, and macOS, Windows, and Linux releases. Distributed under AGPL-3.0. Do not imply query performance from unit-test timings or invent user counts and testimonials.

## Evidence on Hand

Actual application screenshots in `public/screenshots`. Release v0.1.0 has Linux x86_64, macOS Universal, and Windows x64 ZIP assets, verified against the public GitHub release API during this redesign.

## Scope

The user requested a modern, distinctive landing page with good UX in `website/`. No desktop application redesign, deployment, or brand commitments beyond the existing product name were requested.
