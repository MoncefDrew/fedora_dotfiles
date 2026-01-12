# Fedora Dotfiles (Sway Setup)

Hello, this is my new dotfiles repository for the GOAT distro: **Fedora**.

I started with **GNOME**, but ended up setting up **Sway**. It is far less resource-consuming and extremely minimalist.

## Installation

### 1. Clone the repository

```bash
git clone https://github.com/MoncefDrew/fedora-dotfiles
```

### 2. Backup your existing dotfiles

It is highly recommended to back up your current configuration before proceeding.

```bash
mv ~/.config/sway ~/.config/sway_backup
mv ~/.config/nvim ~/.config/nvim_backup
mv ~/.config/neofetch ~/.config/neofetch_backup
mv ~/.config/waybar ~/.config/waybar_backup
mv ~/.config/kitty ~/.config/kitty_backup
```

### 3. Copy the dotfiles

Copy the contents of this repository into your `~/.config` directory.

### 4. Rename Neovim configuration

```bash
mv nvim-config ~/.config/nvim
```

### 5. Reload configuration

Reload Sway:

```
Win + Ctrl + R
```

Restart Waybar:

```bash
pkill waybar && waybar &
```

## Notes

- These dotfiles are tailored for Fedora.
- Make sure required dependencies (sway, waybar, kitty, neovim, neofetch, etc.) are installed.

## Enjoy

Enjoy the rice.

