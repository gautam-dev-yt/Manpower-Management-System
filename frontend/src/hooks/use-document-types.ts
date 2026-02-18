import { useState, useEffect } from 'react';
import { api } from '@/lib/api';
import type { AdminDocumentType } from '@/types';
import { DOC_TYPES, DOC_TYPE_CONFIG, MANDATORY_DOC_TYPES } from '@/lib/constants';
import type { DocTypeConfig } from '@/lib/constants';

/**
 * Fetches document types from the DB.
 * Falls back to hardcoded constants if the API is unavailable.
 */
export function useDocumentTypes() {
    const [types, setTypes] = useState<AdminDocumentType[]>([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        api.documentTypes.list()
            .then(res => setTypes(res.data))
            .catch(() => {
                // Fallback to hardcoded constants
                setTypes(convertConstantsToDocumentTypes());
            })
            .finally(() => setLoading(false));
    }, []);

    return { types, loading };
}

/** Convert hardcoded constants to the AdminDocumentType shape (for fallback). */
function convertConstantsToDocumentTypes(): AdminDocumentType[] {
    return Object.entries(DOC_TYPES).map(([key, displayName], index) => {
        const config: DocTypeConfig = DOC_TYPE_CONFIG[key] || DOC_TYPE_CONFIG.other;
        return {
            id: key,
            docType: key,
            displayName,
            isMandatory: (MANDATORY_DOC_TYPES as readonly string[]).includes(key),
            hasExpiry: config.hasExpiry,
            numberLabel: config.numberLabel,
            numberPlaceholder: config.numberPlaceholder,
            expiryLabel: config.expiryLabel,
            sortOrder: (index + 1) * 10,
            metadataFields: config.metadataFields,
            isSystem: true,
            isActive: true,
            createdAt: '',
            updatedAt: '',
        };
    });
}
