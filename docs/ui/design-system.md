{
  "colorScheme": "light",
  "fonts": [
    {
      "family": "GeistSans",
      "count": 36
    },
    {
      "family": "GeistSans Fallback",
      "count": 36
    },
    {
      "family": "Inter",
      "count": 1
    },
    {
      "family": "Inter Fallback",
      "count": 1
    },
    {
      "family": "Peano",
      "count": 1
    },
    {
      "family": "Liberation Sans",
      "count": 1
    },
    {
      "family": "Arial",
      "count": 1
    }
  ],
  "colors": {
    "primary": "#DBEAFE",
    "accent": "#EFF6FF",
    "background": "#FFFFFF",
    "textPrimary": "#000000",
    "link": "#EFF6FF"
  },
  "typography": {
    "fontFamilies": {
      "primary": "Inter",
      "heading": "GeistSans"
    },
    "fontStacks": {
      "body": [
        "Inter",
        "Inter Fallback",
        "Peano",
        "Liberation Sans",
        "Arial"
      ],
      "heading": [
        "GeistSans",
        "GeistSans Fallback"
      ],
      "paragraph": [
        "GeistSans",
        "GeistSans Fallback"
      ]
    },
    "fontSizes": {
      "h1": "80px",
      "h2": "30px",
      "body": "30px"
    }
  },
  "spacing": {
    "baseUnit": 4,
    "borderRadius": "6px"
  },
  "components": {},
  "images": {
    "logo": "https://storage.googleapis.com/corca-public/about/corca.svg",
    "favicon": "https://corca.app/favicon-light.ico",
    "ogImage": "https://corca.app/corca-og.png"
  },
  "__framework_hints": [],
  "__llm_logo_reasoning": {
    "selectedIndex": 0,
    "reasoning": "Strong indicators: in header, no link (brand logos usually link to homepage, penalty), header location, visible, top-left position, highest logo on page, alt exactly matches brand name \"Corca\", src contains brand name \"Corca\", reasonable size. Score: 98 (clear winner by 98 points)",
    "confidence": 0.9
  },
  "__llm_button_reasoning": {
    "primary": {
      "index": -1,
      "text": "N/A",
      "reasoning": "LLM failed"
    },
    "secondary": {
      "index": -1,
      "text": "N/A",
      "reasoning": "LLM failed"
    },
    "confidence": 0
  },
  "confidence": {
    "buttons": 0,
    "colors": 0,
    "overall": 0
  }
}



{
  "type": "object",
  "description": "Comprehensive design system data from corca.app for shadcn/ui adaptation",
  "properties": {
    "corca_color_palettes": {
      "type": "object",
      "description": "HEX color palettes from corca.app",
      "properties": {
        "backgrounds": {
          "type": "array",
          "description": "Background color palette",
          "items": {
            "type": "string",
            "description": "HEX color code"
          }
        },
        "foregrounds": {
          "type": "array",
          "description": "Foreground color palette",
          "items": {
            "type": "string",
            "description": "HEX color code"
          }
        },
        "primary": {
          "type": "array",
          "description": "Primary color palette",
          "items": {
            "type": "string",
            "description": "HEX color code"
          }
        },
        "accents": {
          "type": "array",
          "description": "Accent color palette",
          "items": {
            "type": "string",
            "description": "HEX color code"
          }
        }
      }
    },
    "corca_border_radius": {
      "type": "array",
      "description": "Border radius values from corca.app",
      "items": {
        "type": "string",
        "description": "Border radius value (e.g., '0.25rem')"
      }
    },
    "corca_typography_scales": {
      "type": "object",
      "description": "Typography scales from corca.app",
      "properties": {
        "font_sizes": {
          "type": "array",
          "description": "Font size scale",
          "items": {
            "type": "string",
            "description": "Font size value (e.g., '1rem')"
          }
        },
        "font_weights": {
          "type": "array",
          "description": "Font weight scale",
          "items": {
            "type": "string",
            "description": "Font weight value (e.g., '500')"
          }
        },
        "line_heights": {
          "type": "array",
          "description": "Line height scale",
          "items": {
            "type": "string",
            "description": "Line height value (e.g., '1.5')"
          }
        }
      }
    },
    "corca_spacing_system": {
      "type": "array",
      "description": "Spacing system from corca.app",
      "items": {
        "type": "string",
        "description": "Spacing value (e.g., '1rem')"
      }
    },
    "corca_component_styling": {
      "type": "object",
      "description": "Specific component styling from corca.app",
      "properties": {
        "buttons": {
          "type": "object",
          "description": "Button styling",
          "properties": {
            "background_color": {
              "type": "string",
              "description": "Button background color (HEX)"
            },
            "text_color": {
              "type": "string",
              "description": "Button text color (HEX)"
            },
            "border_radius": {
              "type": "string",
              "description": "Button border radius"
            },
            "font_size": {
              "type": "string",
              "description": "Button font size"
            },
            "font_weight": {
              "type": "string",
              "description": "Button font weight"
            }
          }
        },
        "inputs": {
          "type": "object",
          "description": "Input styling",
          "properties": {
            "background_color": {
              "type": "string",
              "description": "Input background color (HEX)"
            },
            "text_color": {
              "type": "string",
              "description": "Input text color (HEX)"
            },
            "border_color": {
              "type": "string",
              "description": "Input border color (HEX)"
            },
            "border_radius": {
              "type": "string",
              "description": "Input border radius"
            },
            "font_size": {
              "type": "string",
              "description": "Input font size"
            }
          }
        }
      }
    }
  }
}
{
  "corca_color_palettes": {
    "backgrounds": [
      {
        "value": "#FFFFFF",
        "value_citation": "https://corca.app/"
      },
      {
        "value": "#F3F4F6",
        "value_citation": "https://corca.app/"
      },
      {
        "value": "#EFF6FF",
        "value_citation": "https://corca.app/"
      },
      {
        "value": "#DBEAFE",
        "value_citation": "https://corca.app/"
      },
      {
        "value": "#BFDBFE",
        "value_citation": "https://corca.app/"
      },
      {
        "value": "#1E40AF",
        "value_citation": "https://corca.app/"
      }
    ],
    "foregrounds": [
      {
        "value": "#000000",
        "value_citation": "https://corca.app/"
      },
      {
        "value": "#D1D5DB",
        "value_citation": "https://corca.app/"
      },
      {
        "value": "#9CA3AF",
        "value_citation": "https://corca.app/"
      },
      {
        "value": "#3B82F6",
        "value_citation": "https://corca.app/"
      },
      {
        "value": "#2563EB",
        "value_citation": "https://corca.app/"
      }
    ],
    "primary": [
      {
        "value": "#EFF6FF",
        "value_citation": "https://corca.app/"
      },
      {
        "value": "#DBEAFE",
        "value_citation": "https://corca.app/"
      },
      {
        "value": "#BFDBFE",
        "value_citation": "https://corca.app/"
      },
      {
        "value": "#3B82F6",
        "value_citation": "https://corca.app/"
      },
      {
        "value": "#2563EB",
        "value_citation": "https://corca.app/"
      },
      {
        "value": "#1E40AF",
        "value_citation": "https://corca.app/"
      }
    ],
    "accents": [
      {
        "value": "#E5E7EB",
        "value_citation": "https://corca.app/"
      }
    ]
  },
  "corca_border_radius": [
    {
      "value": "10px",
      "value_citation": "https://corca.app/"
    },
    {
      "value": "20px",
      "value_citation": "https://corca.app/"
    }
  ],
  "corca_typography_scales": {
    "font_sizes": [
      {
        "value": "16px",
        "value_citation": "https://corca.app/"
      },
      {
        "value": "18px",
        "value_citation": "https://corca.app/"
      },
      {
        "value": "28px",
        "value_citation": "https://corca.app/"
      },
      {
        "value": "30px",
        "value_citation": "https://corca.app/"
      },
      {
        "value": "34px",
        "value_citation": "https://corca.app/"
      },
      {
        "value": "46px",
        "value_citation": "https://corca.app/"
      },
      {
        "value": "60px",
        "value_citation": "https://corca.app/"
      },
      {
        "value": "80px",
        "value_citation": "https://corca.app/"
      }
    ],
    "font_weights": [
      {
        "value": "400",
        "value_citation": "https://corca.app/"
      },
      {
        "value": "500",
        "value_citation": "https://corca.app/"
      },
      {
        "value": "600",
        "value_citation": "https://corca.app/"
      }
    ],
    "line_heights": [
      {
        "value": "23.4px",
        "value_citation": "https://corca.app/"
      },
      {
        "value": "45px",
        "value_citation": "https://corca.app/"
      },
      {
        "value": "54px",
        "value_citation": "https://corca.app/"
      },
      {
        "value": "72px",
        "value_citation": "https://corca.app/"
      },
      {
        "value": "30px",
        "value_citation": "https://corca.app/"
      }
    ]
  },
  "corca_spacing_system": [
    {
      "value": "20px",
      "value_citation": "https://corca.app/"
    },
    {
      "value": "43px",
      "value_citation": "https://corca.app/"
    },
    {
      "value": "10px",
      "value_citation": "https://corca.app/"
    },
    {
      "value": "33px",
      "value_citation": "https://corca.app/"
    },
    {
      "value": "8px",
      "value_citation": "https://corca.app/"
    },
    {
      "value": "80px",
      "value_citation": "https://corca.app/"
    },
    {
      "value": "113px",
      "value_citation": "https://corca.app/"
    },
    {
      "value": "40px",
      "value_citation": "https://corca.app/"
    },
    {
      "value": "42px",
      "value_citation": "https://corca.app/"
    },
    {
      "value": "52px",
      "value_citation": "https://corca.app/"
    },
    {
      "value": "60px",
      "value_citation": "https://corca.app/"
    },
    {
      "value": "36px",
      "value_citation": "https://corca.app/"
    },
    {
      "value": "16px",
      "value_citation": "https://corca.app/"
    },
    {
      "value": "27px",
      "value_citation": "https://corca.app/"
    },
    {
      "value": "30px",
      "value_citation": "https://corca.app/"
    },
    {
      "value": "11px",
      "value_citation": "https://corca.app/"
    },
    {
      "value": "14px",
      "value_citation": "https://corca.app/"
    },
    {
      "value": "85px",
      "value_citation": "https://corca.app/"
    },
    {
      "value": "138px",
      "value_citation": "https://corca.app/"
    },
    {
      "value": "13px",
      "value_citation": "https://corca.app/"
    },
    {
      "value": "12px",
      "value_citation": "https://corca.app/"
    },
    {
      "value": "50px",
      "value_citation": "https://corca.app/"
    },
    {
      "value": "9px",
      "value_citation": "https://corca.app/"
    },
    {
      "value": "1px",
      "value_citation": "https://corca.app/"
    },
    {
      "value": "61px",
      "value_citation": "https://corca.app/"
    }
  ],
  "corca_component_styling": {
    "buttons": {
      "background_color": "#DBEAFE",
      "background_color_citation": "https://corca.app/",
      "text_color": "#3B82F6",
      "text_color_citation": "https://corca.app/",
      "border_radius": "10px",
      "border_radius_citation": "https://corca.app/",
      "font_size": "18px",
      "font_size_citation": "https://corca.app/",
      "font_weight": "500",
      "font_weight_citation": "https://corca.app/"
    },
    "inputs": {
      "background_color": null,
      "background_color_citation": null,
      "text_color": null,
      "text_color_citation": null,
      "border_color": null,
      "border_color_citation": null,
      "border_radius": null,
      "border_radius_citation": null,
      "font_size": null,
      "font_size_citation": null
    }
  }
}