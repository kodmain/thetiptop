package mail

import (
	"bytes"
	"fmt"
	"mime"
	"mime/multipart"
	"net/textproto"
	"strings"

	"github.com/kodmain/thetiptop/api/internal/infrastructure/observability/logger"
)

const (
	from                     = "From"
	to                       = "To"
	cc                       = "Cc"
	bcc                      = "Bcc"
	subject                  = "Subject"
	mimeVersion              = "MIME-Version"
	contentType              = "Content-Type"
	contentTransfertEncoding = "Content-Transfer-Encoding"
)

// Mail Représente un e-mail avec toutes les informations nécessaires pour l'envoi.
//
// Fields:
// - From: string L'adresse e-mail de l'expéditeur.
// - To: []string Liste des adresses e-mail des destinataires.
// - Cc: []string Liste des adresses e-mail en copie.
// - Bcc: []string Liste des adresses e-mail en copie cachée.
// - Subject: string Le sujet de l'e-mail.
// - Text: []byte Le contenu en texte brut de l'e-mail.
// - Html: []byte Le contenu en HTML de l'e-mail.
// - Attachments: map[string][]byte Les pièces jointes avec leur nom comme clé.
//
// Methods:
// - IsValid: Vérifie si l'e-mail est valide pour l'envoi.
// - ToSMTP: Construit la représentation SMTP de l'e-mail.
type Mail struct {
	To          []string
	Cc          []string
	Bcc         []string
	Subject     string
	Text        []byte
	Html        []byte
	Attachments map[string][]byte
}

// IsValid Vérifie si l'e-mail a suffisamment d'informations pour être envoyé.
//
// Returns:
// - bool: Vrai si l'e-mail est valide, faux sinon.
func (m *Mail) IsValid() bool {
	// Validation de base: vérifier From, To et au moins un contenu (Text ou Html)
	return len(m.To) > 0 && (len(m.Text) > 0 || len(m.Html) > 0)
}

// ToSMTP Construit la représentation SMTP de l'e-mail pour l'envoi.
//
// Returns:
// - []byte: La représentation SMTP de l'e-mail.
// - error: Une erreur si la construction échoue.
func (m *Mail) Prepare(service ServiceInterface) ([]byte, []string, error) {
	// Création d'un buffer pour construire le message
	var msg bytes.Buffer
	logger.Warn(service)
	fromHeader := service.From()
	if service.Expeditor() != "" {
		// Si FromName contient des caractères non-ASCII, il doit être encodé.
		// Sinon, il peut être utilisé tel quel.
		encodedFromName := mime.QEncoding.Encode("UTF-8", service.Expeditor())
		fromHeader = fmt.Sprintf("%s <%s>", encodedFromName, service.From())
	}

	// Header de base
	logger.Warn(fromHeader)
	header := make(map[string]string)
	header[from] = fromHeader
	header[to] = strings.Join(m.To, ", ")
	header[subject] = m.Subject
	header[mimeVersion] = "1.0"

	// Gestion de CC et BCC
	if len(m.Cc) > 0 {
		header[cc] = strings.Join(m.Cc, ", ")
	}

	// Début de la composition MIME
	writer := multipart.NewWriter(&msg)
	boundary := writer.Boundary()
	header[contentType] = "multipart/alternative; boundary=" + boundary + "; charset=UTF-8"
	header[contentTransfertEncoding] = "quoted-printable"

	for key, value := range header {
		msg.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
	}

	msg.WriteString("\r\n")

	if len(m.Text) > 0 {
		if err := m.addContent(writer, "text/plain", m.Text); err != nil {
			return nil, nil, err
		}
	}

	if len(m.Html) > 0 {
		if err := m.addContent(writer, "text/html", m.Html); err != nil {
			return nil, nil, err
		}
	}

	// Ajout des pièces jointes
	for name, data := range m.Attachments {
		if err := m.addAttachment(writer, name, data); err != nil {
			return nil, nil, err
		}
	}

	// Clôture de la composition MIME
	writer.Close()

	// Compilation de tous les destinataires (To, Cc, Bcc)
	allRecipients := append(m.To, m.Cc...)
	allRecipients = append(allRecipients, m.Bcc...)

	return msg.Bytes(), allRecipients, nil
}

// addContent Ajoute du contenu (texte ou HTML) à la composition MIME.
//
// Parameters:
// - writer: *multipart.Writer Le writer pour la composition MIME.
// - contentType: string Le type de contenu (text/plain ou text/html).
// - content: []byte Le contenu à ajouter.
//
// Returns:
// - error: Une erreur si l'ajout échoue.
func (m *Mail) addContent(writer *multipart.Writer, ct string, content []byte) error {
	if len(content) == 0 {
		return nil
	}

	header := make(textproto.MIMEHeader)
	header.Set(contentType, ct+"; charset=UTF-8")
	header.Set("Content-Transfer-Encoding", "quoted-printable")

	part, err := writer.CreatePart(header)
	if err != nil {
		return err
	}

	_, err = part.Write(content)
	return err
}

// addAttachment Ajoute une pièce jointe à la composition MIME.
//
// Parameters:
// - writer: *multipart.Writer Le writer pour la composition MIME.
// - name: string Le nom de la pièce jointe.
// - data: []byte Les données de la pièce jointe.
//
// Returns:
// - error: Une erreur si l'ajout échoue.
func (m *Mail) addAttachment(writer *multipart.Writer, name string, data []byte) error {
	if len(data) == 0 {
		return nil
	}

	header := make(textproto.MIMEHeader)
	header.Set(contentType, "application/octet-stream")
	header.Set("Content-Transfer-Encoding", "base64")
	header.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", name))

	part, err := writer.CreatePart(header)
	if err != nil {
		return err
	}

	_, err = part.Write(data)
	return err
}
