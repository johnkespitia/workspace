import { colors } from '../tokens/colors';
import { spacing } from '../tokens/spacing';
import { typography } from '../tokens/typography';

export const darkTheme = {
  colors: {
    background: colors.dark.neutral[50],
    surface: colors.dark.neutral[100],
    surfaceElevated: colors.dark.neutral[200],
    text: {
      primary: colors.dark.neutral[900],
      secondary: colors.dark.neutral[600],
      tertiary: colors.dark.neutral[400],
    },
    border: colors.dark.neutral[300],
    primary: colors.dark.primary[600],
    primaryHover: colors.dark.primary[700],
    secondary: colors.dark.secondary[600],
    success: colors.dark.success[600],
    error: colors.dark.error[600],
    warning: colors.dark.warning[600],
  },
  spacing,
  typography,
} as const;

export type DarkTheme = typeof darkTheme;
