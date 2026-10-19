---
name: DBStudio Lite website
description: A clear desktop workspace, presented on paper with forest accents.
colors:
  paper: "#fafbf9"
  ink: "#17231e"
  muted: "#58635d"
  line: "#dce2dc"
  accent: "#285b44"
  accent-hover: "#1b4532"
  mint: "#e7eee5"
  preview: "#dce8d5"
  white: "#ffffff"
  download-muted: "#b6ccb7"
  download-line: "#4c6152"
typography:
  display:
    fontFamily: "Arial, Helvetica, sans-serif"
    fontSize: "clamp(48px, 5.8vw, 76px)"
    fontWeight: 500
    lineHeight: 1.05
    letterSpacing: "-.04em"
  headline:
    fontFamily: "Arial, Helvetica, sans-serif"
    fontSize: "clamp(30px, 3.4vw, 44px)"
    fontWeight: 500
    lineHeight: 1.15
    letterSpacing: "-.035em"
  body:
    fontFamily: "Arial, Helvetica, sans-serif"
    fontSize: "14px"
    lineHeight: 1.7
  button:
    fontFamily: "Arial, Helvetica, sans-serif"
    fontSize: "14px"
    fontWeight: 600
rounded:
  control: "6px"
  button: "7px"
spacing:
  compact: "8px"
  inline: "14px"
  group: "24px"
  section-detail: "30px"
components:
  button-primary:
    backgroundColor: "{colors.accent}"
    textColor: "{colors.white}"
    typography: "{typography.button}"
    rounded: "{rounded.button}"
    padding: "0 23px"
  button-primary-hover:
    backgroundColor: "{colors.accent-hover}"
  button-ink:
    backgroundColor: "{colors.ink}"
    textColor: "{colors.white}"
    rounded: "{rounded.button}"
  button-ink-hover:
    backgroundColor: "{colors.accent}"
  preview-tab-selected:
    backgroundColor: "{colors.paper}"
    textColor: "{colors.ink}"
    rounded: "{rounded.control}"
    padding: "11px 15px"
---

# Design System: DBStudio Lite website

## Overview

**Creative North Star: "A clearer workspace"**

The implemented website places the real desktop application in a quiet editorial setting. Paper surfaces, forest text, and sage framing give the screenshots room to explain the product. This records the website implementation, not the desktop application's theme.

The layout alternates broad product evidence with compact supporting details. Familiar system typography keeps the page lightweight; selective color and generous spacing establish hierarchy.

**Key Characteristics:**
- Paper surfaces and forest accents.
- Real application screenshots with user-controlled views.
- Asymmetric desktop composition that stacks on mobile.
- Direct platform downloads and visible keyboard focus.

## Colors

A restrained green family links navigation, actions, product framing, and the dark download area.

### Primary
- **Forest accent:** Primary download actions, highlighted headline text, small icons, and the default focus outline. The deeper hover token supplies button feedback.
- **Mint and sage:** Mint supports the editor feature; sage frames the screenshot showcase.

### Neutral
- **Paper:** Page canvas and selected preview tabs.
- **Forest ink:** Primary text and the compact header action.
- **Muted:** Supporting descriptions and secondary labels.
- **Line:** Hairline section and detail separators.
- **White:** Text on filled actions.
- **Download muted and line:** Supporting labels and dividers on the dark download surface.

## Typography

Display and body use Arial with Helvetica and generic sans-serif fallbacks. Code samples use monospace. No remote font is required.

Display and section headlines use moderate weight, tight tracking, and balanced wrapping. Body text stays looser for scanning. The scale is role-based rather than a fixed mathematical ratio. Feature headings sit between body copy and section headings; metadata is visually subordinate and should remain short.

**The One Family Rule.** Use the existing system sans-serif stack for website prose and controls; reserve monospace for code and database types.

## Layout

The desktop container is capped at 1120px with 40px side gutters. The hero pairs a larger headline column with a smaller description/action column (1.35:1). Features pair the editor example with stacked supporting details (1.2:1). Download options form three equal columns.

At 1050px and below, desktop gaps narrow. At 760px and below, side gutters become 20px, hero and features stack, downloads become full-width rows, and the main navigation links hide while the brand and download action remain. Below 380px gutters become 16px. The header is in normal document flow.

Section rhythm uses broad vertical gaps around the hero, features, and downloads; compact spacing groups labels with actions. Borders separate details without enclosing every item in a card. Anchor scrolling reserves 90px above targets.

## Elevation & Depth

Most surfaces are flat, separated by tone and thin borders. The application window alone has a diffuse resting shadow (`0 18px 40px -20px #25372180`) to distinguish the real product preview from the page. Buttons rise by 1px on hover.

**The Product Depth Rule.** Keep the existing ambient shadow associated with the application preview; supporting sections rely on tone and dividers.

## Shapes

Controls use gently rounded corners. The application frame clips its screenshot within a 10px radius; the editor feature uses the existing 12px radius variable. Those container dimensions are local treatments, not a generalized radius scale. Small circular dots mark status and window chrome. Most dividers are 1px and square-ended.

## Components

### Buttons

Filled actions are compact, with icon and label grouped together. Main buttons have a 50px minimum height and 23px horizontal padding; the header variant uses 40px and 17px. Hover changes the fill and raises the control slightly over 160ms. Focus is a visible 2px accent outline with a 5px offset.

The hero detects desktop operating systems after mount and points to that platform's release. Mobile, unknown platforms, and the initial render lead to the download section. An adjacent link always exposes platform choice.

### Navigation

The wordmark sits opposite a download action. Desktop links connect the workspace preview, features, and repository. Mobile preserves the direct download action. A keyboard-visible skip link reaches main content; text links underline on hover where appropriate.

### Workspace preview

Four tabs select real screenshots and matching captions. The selected tab has a paper fill; unselected tabs share the sage surface. Roving tab focus, Left/Right arrows, Home, and End support keyboard selection. Each tab owns an associated panel, with inactive panels hidden.

Screenshots link to full-size assets. The full-size affordance appears on hover and keyboard focus, and remains visible on mobile. Data and SQL screenshots crop within the frame; connection and profile views fit fully. Desktop framing uses a 2.35 aspect ratio; mobile uses 1.55.

The preview enters with a short upward movement on widths of at least 900px. Reduced-motion preference removes animations, transitions, and smooth scrolling.

### Feature containers

The editor feature uses a mint background and embedded dark, monospace SQL example. Supporting features remain open on paper, with a divider and compact schema sample. The example labels distinguish illustrative code and metadata from interactive product output.

### Platform downloads

A dark forest section groups three links with platform, architecture, icon, and download indicator. Dividers separate platforms; hover changes the row background. Focus uses a light outline inset by 4px. Mobile presents vertical rows. A separate release link provides access to notes and additional assets.

## Do's and Don'ts

### Do:
- **Do** preserve real product screenshots and full-size access.
- **Do** keep platform choice available beside the detected download action.
- **Do** preserve keyboard tab controls, visible focus, and reduced-motion behavior.
- **Do** use the existing system font stack and compact supporting labels.

### Don't:
- **Don't** replace the screenshot preview with invented product UI.
- **Don't** make hover the only way to discover full-size previews on touch devices.
- **Don't** promote local container dimensions into a universal component scale.
