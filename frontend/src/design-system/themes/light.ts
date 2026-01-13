import { colors } from '../tokens/colors';
import { spacing } from '../tokens/spacing';
import { typography } from '../tokens/typography';

export const lightTheme = {
  colors: {
    background: colors.light.neutral[50],
    surface: colors.light.neutral[100],
    surfaceElevated: '#ffffff',
    text: {
      primary: colors.light.neutral[900],
      secondary: colors.light.neutral[600],
      tertiary: colors.light.neutral[400],
    },
    border: colors.light.neutral[200],
    primary: colors.light.primary[600],
    primaryHover: colors.light.primary[700],
    secondary: colors.light.secondary[600],
    success: colors.light.success[600],
    error: colors.light.error[600],
    warning: colors.light.warning[600],
  },
  spacing,
  typography,
} as const;

export type LightTheme = typeof lightTheme;
