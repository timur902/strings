package unpack

import (
	"strings"

	"github.com/timur902/stack_queue/internal/repository"
)

func NewProvider(repo repository.Repository) *Provider {

	return &Provider{
		repo: repo,
	}
}

type Provider struct {
	repo repository.Repository
}

func (p *Provider) Pack(s string) string {

	if len(s) == 0 {
		return ""
	}

	rs := []rune(s)

	var builder strings.Builder

	count := 1

	for i := 1; i <= len(rs); i++ {

		if i < len(rs) && rs[i] == rs[i-1] {
			count++
			continue
		}

		builder.WriteRune(rs[i-1])

		if count > 1 {
			builder.WriteString(string(rune('0' + count)))
		}

		count = 1
	}

	p.repo.InsertString(builder.String())



	return builder.String()
}

// Получить по id вернуть результат если нету, то ошибку.

func (p *Provider) Unpack(s string) (string, error) {
	rs := []rune(s)
	if len(rs) == 0 {
		return "", nil
	}
	var out []rune
	var prev rune
	hasPrev := false
	prevWasDigit := false
	escaped := false
	for i := 0; i < len(rs); i++ {
		r := rs[i]
		if escaped {
			if (r < '0' || r > '9') && r != '\\' {
				return "", ErrInvalidString
			}
			out = append(out, r)
			prev = r
			hasPrev = true
			prevWasDigit = false
			escaped = false
			continue
		}
		if r == '\\' {
			escaped = true
			continue
		}
		if r >= '0' && r <= '9' {

			if !hasPrev {
				return "", ErrInvalidString
			}

			if prevWasDigit {
				return "", ErrInvalidString
			}

			count := int(r - '0')

			if count == 0 {
				out = out[:len(out)-1]
				prevWasDigit = true
				continue
			}

			for j := 0; j < count-1; j++ {
				out = append(out, prev)
			}

			prevWasDigit = true
			continue
		}

		out = append(out, r)
		prev = r
		hasPrev = true
		prevWasDigit = false
	}
	if escaped {
		return "", ErrInvalidString
	}
	return string(out), nil
}
