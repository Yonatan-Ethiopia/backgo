package main

import "backgo/internal/models"

type templateData struct{
    Box *models.RecRow
    Boxes []*models.RecRow
}
