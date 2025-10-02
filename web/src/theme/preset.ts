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
})

export default SnagglePreset
