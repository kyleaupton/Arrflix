import { definePreset } from '@primeuix/themes'
import Aura from '@primeuix/themes/aura'

// Snaggle theme preset based on Aura. Customize tokens here.
const SnagglePreset = definePreset(Aura, {
  primitive: {
    borderRadius: {
      sm: '8px',
      md: '10px',
      lg: '12px',
      xl: '14px',
    },
  },
  semantic: {
    transitionDuration: '150ms',
    focusRing: {
      width: '2px',
      style: 'solid',
      color: '#7c3aed',
      offset: '2px',
      shadow: '0 0 0 2px color-mix(in oklab, #7c3aed, transparent 70%)',
    },
    shadow: {
      1: 'inset 0 1px 2px #ffffff30, 0 1px 2px #00000030, 0 2px 4px #00000015', // --shadow-sm
      2: 'inset 0 1px 2px #ffffff50, 0 2px 4px #00000030, 0 4px 8px #00000015', // --shadow-md
      3: 'inset 0 1px 2px #ffffff70, 0 4px 6px #00000030, 0 6px 10px #00000015', // --shadow-lg
    },
    colorScheme: {
      light: {
        primary: {
          color: '#7c3aed',
          contrastColor: '#ffffff',
          hoverColor: '#6d28d9',
          activeColor: '#5b21b6',
        },
      },
      dark: {
        primary: {
          color: '#8b5cf6',
          contrastColor: '#0b1020',
          hoverColor: '#a78bfa',
          activeColor: '#7c3aed',
        },
      },
    },
  },
  components: {
    inputtext: {
      root: {
        shadow: '{shadow.1}',
        borderColor: 'transparent',
      },
    },
    menu: {
      root: {
        shadow: '{shadow.1}',
        borderColor: 'transparent',
      },
    },
    panelmenu: {
      root: {
        background: 'transparent',
        borderColor: 'transparent',
      },
    },
    datatable: {
      root: {
        // shadow: '{shadow.1}',
      },
      header: {
        padding: '0.5rem 1rem',
      },
    },
    card: {
      root: {
        shadow: '{shadow.1}',
      },
    },
  },
})

// Export CSS custom properties for layout dimensions
// These can be used throughout the app for consistent spacing
export const layoutTokens = {
  sidebarWidth: '280px',
  headerHeight: '64px',
  layoutGap: '1rem',
  layoutPadding: '1rem',
  layoutPaddingMobile: '0.75rem',
} as const

export default SnagglePreset
